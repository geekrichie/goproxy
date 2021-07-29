package mux_net

import (
	"encoding/binary"
	"errors"
	"fmt"
	"goproxy/common"
	"goproxy/log"
	"goproxy/mux/mux_link"
	"goproxy/mux/mux_msg"
	"io"
	"net"
	"sync/atomic"
)

var errConnectFail = errors.New("connect to the main server failed")
var errHandshake  = errors.New("handshake with server failure")
var errSecretkey = errors.New("key error")



type Connection struct {
	conn     net.Conn
	ConnType int
	Id       int
	Plexer   *mux_link.MultiPlexer
}


func NewConn(conn net.Conn) Connection {
	return Connection{
		conn : conn,
	}
}

func (c *Connection) SetConnType(connType int) {
	c.ConnType = connType
}

func (c *Connection) SendLinkInfo(targetaddr string) {
	//这里1个字节的类型标识，4个字节的长度，后面接具体的连接信息
	c.SendMsg(mux_msg.MSG_LINK_INFO, targetaddr)
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
	if err != nil || string(msgConnect) != "connected" {
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
	_, err = c.Write([]byte("connected"))
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
	var l int32
	err := binary.Read(c, binary.LittleEndian, &l)
	return int(l), err
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

func (c *Connection) SendMsg(msgType uint8, message string) {
	//这里的发送消息格式是
	//-------------------------------
	//| msgType | msglen |   msg   |
	//-------------------------------
	//| 1byte   | 4byte  | msglen值 |
	//-------------------------------
	msgInfo := mux_msg.SyncMsgInfoPool.Get().(*mux_msg.MsgInfo)
	msgInfo.SetMessage(msgType, message)
	packedMessage := msgInfo.Pack()
	_, err := c.Write(packedMessage)
	if err != nil {
		log.Error(err.Error())
	}
}

func Unpack(conn Connection, msgType uint8) {
	var buf = make([]byte, 4)
	io.ReadFull(&conn, buf)
	var connId uint32
	connId = binary.LittleEndian.Uint32(buf)
	linkConn := conn.Plexer.GetConnById(int(connId))
	if msgType == mux_msg.MSG_CLOSE_CONN {
		if linkConn.GetConnId() != 0 {
			linkConn.Close()
		}
		return
	}
	io.ReadFull(&conn, buf)
	var messagelen uint32
	messagelen = binary.LittleEndian.Uint32(buf)
	//log.Infof("New messagelen : %d", messagelen)
	var message = make([]byte, messagelen)
	io.ReadFull(&conn, message)
	linkConn.ReceiveWindowWrite(message)
	return
}

/**
 kafka-go balance.go 中的轮询实现
 */

type Target struct {
	offset uint32
	TargetAddrs []string //tcp代理访问的目标地址
}



func (t *Target) GetRandomAddr() string{
	offset := atomic.AddUint32(&t.offset, 1) - 1
	return t.TargetAddrs[offset%uint32(len(t.TargetAddrs))]

}
