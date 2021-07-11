package mux

import (
	"encoding/binary"
	"net"
	"sync"
	"sync/atomic"
)


const (
	TYPE_LINK_INFO = iota
)

type Connection struct {
	conn net.Conn
	Target  Target
}


func NewConn(conn net.Conn) Connection{
	return Connection{
		conn : conn,
	}
}

func (c *Connection) SendLinkInfo(targetaddr string)error {
	//这里1个字节的类型标识，4个字节的长度，后面接具体的连接信息
	return c.Pack(TYPE_LINK_INFO, targetaddr)
}

func (c *Connection) Pack(packetType uint8, info string) (err error) {
	err = binary.Write(c.conn, binary.LittleEndian, packetType)
	if err != nil {
		return
	}
	err = binary.Write(c.conn, binary.LittleEndian, len(info))
	if err != nil {
		return
	}
	err = binary.Write(c.conn, binary.LittleEndian, info)
	return
}

type Target struct {
	index int64
	lock sync.Locker
	TargetAddrs []string //tcp代理访问的目标地址
}

func (t *Target) GetRandomAddr() string{
     atomic.StoreInt64(&t.index, t.index+1)
     return t.TargetAddrs[t.index%int64(len(t.TargetAddrs))]
}
