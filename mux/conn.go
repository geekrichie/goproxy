package mux

import (
	"encoding/binary"
	"errors"
	"net"
	"sync"
	"sync/atomic"
)

var errConnectFail = errors.New("connect to the main server failed")

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

func (c *Connection) Read(b []byte) (n int,err error) {
	return c.conn.Read(b)
}

func (c *Connection) Write(b []byte) (n int, err error) {
	return c.conn.Write(b)
}

func (c *Connection) SendHandShake() error{
	var (
		err error
		n  int
	)
	_, err = c.Write([]byte("connect"))
	if err != nil {
		return err
	}
	var buf = make([]byte, 20)
	n, err = c.Read(buf)
	msgConnect := buf[:n]
	if string(msgConnect) != "connect ok" {
		return err
	}

	return nil
}

func (c *Connection) ReceiveHandShake() error{
	var (
		err error
		n  int
	)
	var buf = make([]byte, 20)
	n, err = c.Read(buf)
	msgSend := buf[:n]
	if string(msgSend) != "connect" {
		return err
	}
	_, err = c.Write([]byte("connect ok"))
	return err
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
