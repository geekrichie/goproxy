package mux_link

import (
	"goproxy/mux/mux_queue"
	"net"
	"sync"
	"time"
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
	conns []conn
	L sync.Mutex
	netconn net.Conn
}

type conn struct {
	connId int
	readWindow receiveWindow
	sendWindow sendWindow
}

func NewMultiPlexer() *MultiPlexer {
	return &MultiPlexer{
		connNum: 0,
		conns : make([]conn,0),
	}
}

func (m MultiPlexer) AddConn(conn conn) {
	m.L.Lock()
	defer m.L.Unlock()
	conn.connId = m.connNum
	m.connNum  = m.connNum + 1
	m.conns = append(m.conns, conn)
}




func (c *conn) GetConnId() int{
	return c.connId
}
//
//func (c *conn) Write([]byte) int {
//
//}
//
//func (c *conn) Read([]byte) {
//
//}

type receiveWindow struct {
	bufQueue *mux_queue.LKQueue
	curElem *ListElement
	//当前读取List的偏移量
	off int
	windowSize  int
	maxWindowsize int
}

func NewReceiveWindow() *receiveWindow{
	 return &receiveWindow{
           bufQueue: mux_queue.NewLKQueue(),
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
	rw.windowSize += len(b)
	listelem := syncListPool.Get().(*ListElement)
	listelem.Buf = b
	listelem.L = len(b)
	rw.bufQueue.Enqueue(listelem)

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
	startRead:
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
		s := rw.bufQueue.Dequeue()
		if s == nil{
			time.Sleep(1*time.Millisecond)
			goto startRead
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
	rw.windowSize -= len(b)
	return len(b), nil
}


type sendWindow struct {

}

