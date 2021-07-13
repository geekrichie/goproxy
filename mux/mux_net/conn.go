package mux_net

import (
	"encoding/binary"
	"errors"
	"fmt"
	"goproxy/common"
	"goproxy/log"
	"io"
	"net"
	"sync"
	"sync/atomic"
)

var errConnectFail = errors.New("connect to the main server failed")
var errHandshake  = errors.New("handshake with server failure")
var errSecretkey = errors.New("key error")

const (
	TYPE_LINK_INFO = iota
)

type Connection struct {
	conn     net.Conn
	Target   Target
	ConnType int
	Id       int
}


func NewConn(conn net.Conn) Connection {
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

func (c *Connection) Close() error {
	return c.conn.Close()
}

func (c *Connection) SendMode(mode uint8)(n int, err error) {
	return c.Write([]byte{mode})
}

func (c *Connection) SendHandShake() error{
	var (
		err error
	)
	_, err = c.Write([]byte("connect"))
	if err != nil {
		return err
	}
	msgConnect,err := c.ReadSmallMsg()
	if err != nil || string(msgConnect) != "connect ok" {
		if err != nil{
			return err
		}
		return errHandshake
	}
	log.Infof("connect to the server %s ,received %s\n", c.conn.RemoteAddr().String(), string(msgConnect))
	_, err = c.Write([]byte(common.GetSecretKey()))
	if err != nil{
		return err
	}
	resp, err := c.ReadSmallMsg()
    if err != nil || string(resp) != "right key" {
    	if err != nil {
    		return err
		}
		return errSecretkey
	}
	return nil
}

func (c *Connection) ReadSmallMsg() ([]byte, error){
	var (
		err error
		n  int
	)
	var buf = make([]byte, 20)
	n, err = c.Read(buf)
	return buf[:n], err
}

func (c *Connection) ReceiveHandShake() error{
	var (
		err error
	)
	msgSend , err := c.ReadSmallMsg()
	if err!= nil || string(msgSend) != "connect" {
		if err != nil {
			return err
		}
		return errors.New(fmt.Sprintf("client: %s wrong connect msg", c.conn.RemoteAddr().String()))
	}
	_, err = c.Write([]byte("connect ok"))
	if err != nil {
		return err
	}

	secretkey,err := c.ReadSmallMsg()
	if err != nil || string(secretkey) != common.GetSecretKey() {
		if err != nil {
			return err
		}
		c.Write([]byte("wrong key"))
		return errSecretkey
	}
	_, err = c.Write([]byte("right key"))

	return err
}

func (c *Connection) ReadLenContent() ([]byte, error){
	l, err := c.ReadLen()
	if err != nil {
		return nil,err
	}
	buf, err := c.ReadContent(l)
	if err != nil {
		return nil,err
	}
	return buf, nil
}

func (c *Connection) ReadLen() (int,error){
	var l int
	err := binary.Read(c, binary.LittleEndian, &l)
	return l, err
}

func (c *Connection) ReadMsgType() (uint8,error){
	var l uint8
	err := binary.Read(c, binary.LittleEndian, &l)
	return l, err
}

func (c *Connection) ReadContent(contentSize int)([]byte,error) {
	var buf  = make([]byte,contentSize)
	_,err := io.ReadFull(c, buf)
	return buf, err
}

func (c *Connection) GetConn() net.Conn{
	return c.conn
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