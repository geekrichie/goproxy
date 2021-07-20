package mux_link

import (
	"net"
	"sync"
)

const (
	MainMode uint8 = iota
	TranMode
)

type MultiPlexer struct {
	connNum int
	conns []conn
	L sync.Mutex
	netconn net.Conn
}

type conn struct {
	connId int
	//readQueue Queue
	//writeQueue Queue
}

func NewMultiPlexer() *MultiPlexer {
	return &MultiPlexer{
		connNum: 0,
		conns : make([]conn,10),
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

func (c *conn) Write() int {

}

