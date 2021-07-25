package mux_link

import (
	"goproxy/mux/mux_queue"
	"net"
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
	queueSize int
}

func NewReceiveWindow() *receiveWindow{
	 return &receiveWindow{
           bufQueue: mux_queue.NewLKQueue(),
           queueSize: 0,
	 }
}

func (rw *receiveWindow) Write(b []byte) {
	rw.bufQueue.Enqueue(b)

}


type sendWindow struct {

}

