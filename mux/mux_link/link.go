package mux_link

import (
	"goproxy/log"
	"goproxy/mux/mux_msg"
	"goproxy/mux/mux_queue"
	"net"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"
)

const (
	MainMode uint8 = iota
	TranMode
)

const (
	MaxReceiveWindowSize = 4* 1024*1024
)

type MultiPlexer struct {
	connNum int
	conns  map[int]*conn
	L sync.Mutex
	netconn net.Conn
}

type conn struct {
	connId int
	plexer *MultiPlexer
	receiveWindow receiveWindow
	sendWindow sendWindow
	isClose bool
	once sync.Once
}

func (c *conn) SetDeadline(t time.Time) error {
	return c.plexer.netconn.SetDeadline(t)
}

func (c *conn) SetReadDeadline(t time.Time) error {
	return c.plexer.netconn.SetReadDeadline(t)
}

func (c *conn) SetWriteDeadline(t time.Time) error {
	return c.plexer.netconn.SetWriteDeadline(t)
}

func NewMultiPlexer(netconn net.Conn) *MultiPlexer {
	return &MultiPlexer{
		connNum: 0,
		conns : make(map[int]*conn),
		netconn: netconn,
		L : sync.Mutex{},
	}
}

func (m *MultiPlexer) AddConn(conn *conn) {
	m.L.Lock()
	defer m.L.Unlock()
	m.connNum  = m.connNum + 1
	conn.connId = m.connNum
	conn.plexer = m
	m.conns[conn.connId] = conn
}

func (m *MultiPlexer) RemoveConn(conn *conn) {
	m.L.Lock()
	defer m.L.Unlock()
	if _,ok := m.conns[conn.connId]; ok {
		delete(m.conns,conn.connId)
		m.connNum = m.connNum - 1
		conn.plexer = nil
	}
}

func (m *MultiPlexer)GetConnById(connId int) *conn{
	return m.conns[connId]
}

func (m *MultiPlexer) Write(data []byte)(n int, err error) {
	return m.netconn.Write(data)
}
func (m *MultiPlexer) Read(data []byte)(n int, err error) {
	return m.netconn.Read(data)
}

func NewConn(plexer *MultiPlexer) *conn{
	return &conn{
		receiveWindow: *NewReceiveWindow(),
		sendWindow: sendWindow{plexer: plexer},
	}
}

func (c *conn) SendLinkInfo(targetaddr string) {
	//这里1个字节的类型标识，4个字节的长度，后面接具体的连接信息
	c.SendInfo(mux_msg.MSG_LINK_INFO, targetaddr)
}

func (c *conn) SendInfo(flag uint8, message string) {
	//这里1个字节的类型标识，4个字节的长度，后面接具体的连接信息
	msgConnInfo := mux_msg.SyncMsgConnInfoPool.Get().(*mux_msg.MsgConnInfo)
	msgConnInfo.SetMessage(flag, int32(c.connId), message)
	buf, err := msgConnInfo.Pack()
	if err != nil {
		log.Error(err.Error())
	}
	_, err = c.sendWindow.plexer.Write(buf)
	if err != nil {
		log.Error(err.Error())
	}
}


func Copy(c1 , c2 net.Conn) {
	buf := make([]byte,32*1024)
	go func() {
		for {
			n1, err1 := c1.Read(buf)
			if err1 == nil {
				c2.Write(buf[:n1])
			}else {
				log.Error(err1.Error())
				//发生错误两端关闭
				c1.Close()
				c2.Close()
				break
			}
		}
	}()
	go func() {
		for {
			n2, err2 := c2.Read(buf)
			if err2 == nil {
				c1.Write(buf[:n2])
			} else {
				log.Error(err2.Error())
				//发生错误两端关闭
				c1.Close()
				c2.Close()
				break
			}
		}
	}()
}


func (c *conn) GetConnId() int{
	return c.connId
}

func (c *conn) SetConnId(connId int) {
	c.connId = connId
}

func (c *conn) Write(data []byte) (n int, err error) {
   return  c.sendWindow.Write(data, c.connId)
}

func (c *conn) Read(data []byte) (n int, err error) {
    return c.receiveWindow.Read(data)
}
func (c *conn) Close() error {
	//todo 关闭连接
	//这里会先关闭连接，然后通知对方关闭
	c.once.Do(func() {
		 c.isClose = true
		 c.SendInfo(mux_msg.MSG_CLOSE_CONN,"")
		 c.plexer.RemoveConn(c)
	})
	return nil
}

// LocalAddr returns the local network address.
func (c *conn) LocalAddr() net.Addr {
	return c.plexer.netconn.LocalAddr()
}

// RemoteAddr returns the remote network address.
func (c *conn) RemoteAddr() net.Addr {
	return c.plexer.netconn.RemoteAddr()
}

func (c *conn) ReceiveWindowWrite(message []byte)error {
	_,err := c.receiveWindow.Write(message)
	return err
}

type receiveWindow struct {
	bufQueue *mux_queue.Queue
	curElem *ListElement
	//当前读取List的偏移量
	off int
	windowSize  int
	maxWindowsize int
	plexer MultiPlexer
}
//todo windowSize 不知道怎么为负数了

func NewReceiveWindow() *receiveWindow{
	 return &receiveWindow{
           bufQueue: mux_queue.NewQueue(),
           curElem : new(ListElement),
           windowSize: 0,
           maxWindowsize: MaxReceiveWindowSize,
	 }
}

func (rw *receiveWindow) Write(b []byte)(n int, err error) {
	//如果超出最大窗口大小，那么先阻塞
	start:
	remainsize :=   MaxReceiveWindowSize-rw.windowSize
	if remainsize < len(b) {
		time.Sleep(time.Millisecond * 5)
		goto start
	}
	atomic.StoreInt32((*int32)(unsafe.Pointer(&rw.windowSize)), int32(rw.windowSize+len(b)))
	//log.Info(fmt.Sprintf(" write receive window size : %d", rw.windowSize))
	listelem := syncListPool.Get().(*ListElement)
	listelem.Buf = b
	listelem.L = len(b)
	rw.bufQueue.Push(listelem)

	return len(b), nil
}

type ListElement struct{
	Buf    []byte
	L      int
	isPart bool
}

var syncListPool = sync.Pool{
	New: func()interface{} {
		return new(ListElement)
     },
}

func (rw *receiveWindow) Read(b []byte)(int, error)  {
	var n int
	var s interface{}
	var err error
	off := 0
	if rw.curElem.L > rw.off {
		n = copy(b,rw.curElem.Buf[rw.off:rw.curElem.L])
		off += n
		rw.off += n
	}
	if off == len(b) {
		return n, nil
	}
	for {
		if rw.bufQueue.Len == 0 {
			if off == 0 {
				//如果什么都没读到就阻塞，只要读到了一部分就返回
				time.Sleep(10*time.Millisecond)
			}else {
				break
			}
		}
		s,err = rw.bufQueue.Pop()
		if err != nil {
			continue
		}
		elem := s.(*ListElement)
		//off表示已经读取的长度，如果当前元素的长度大于要读取的长度

		n = copy(b[off:], elem.Buf[:elem.L])
		off += n
		if off >= len(b) {
			rw.curElem = elem
			rw.off = n
			break
		}
	}
	//log.Info(fmt.Sprintf("read before receive window size : %d", rw.windowSize))
	//rw.windowSize -= off
	atomic.StoreInt32((*int32)(unsafe.Pointer(&rw.windowSize)), int32(rw.windowSize-off))
	//log.Info(fmt.Sprintf("read receive window size : %d", rw.windowSize))
	return off, err
}


type sendWindow struct {
	plexer *MultiPlexer
}


func (s *sendWindow) Write(data []byte, connId int)(n int, err error) {
	msgConnInfo := mux_msg.SyncMsgConnInfoPool.Get().(*mux_msg.MsgConnInfo)
	msgConnInfo.SetMessage(mux_msg.MSG_TRAN_INFO, int32(connId),string(data))
	packedData,err := msgConnInfo.Pack()
	if err != nil {
		log.Error(err.Error())
	}
	return s.plexer.netconn.Write(packedData)
}
