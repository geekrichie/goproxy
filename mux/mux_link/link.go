package mux_link

import (
	"goproxy/mux/mux_queue"
	"net"
	"runtime"
	"sync"
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
	windowSize  int
	maxWindowsize int
}

func NewReceiveWindow() *receiveWindow{
	 return &receiveWindow{
           bufQueue: mux_queue.NewLKQueue(),
           windowSize: 0,
           maxWindowsize: MaxReceiveWindowSize,
	 }
}

func (rw *receiveWindow) Write(b []byte)(n int, err error) {
	//如果超出最大窗口大小，那么先阻塞
	if rw.windowSize + len(b) > MaxReceiveWindowSize{
		runtime.Gosched()
	}
	rw.bufQueue.Enqueue(b)
	rw.windowSize += len(b)
	return
}

func (rw *receiveWindow) Read(b []byte)(n int, err error)  {

	return
}


type sendWindow struct {

}

