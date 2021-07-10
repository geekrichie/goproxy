package mux

import (
	"net"
	"sync"
)

type Connection struct {
	conn net.Conn
	Target   Target
}

type Target struct {
	lock sync.Locker
	TargetAddrs []string //tcp代理访问的目标地址
}

func NewConn(conn net.Conn) Connection{
	return Connection{
		conn : conn,
	}
}


func (c *Connection) GetRandomAddr() {


}