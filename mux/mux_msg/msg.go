package mux_msg

import (
	"bytes"
	"encoding/binary"
	"sync"
)

const (
	MSG_PING uint8 = iota
	MSG_LOG_INFO
	MSG_LINK_INFO
	MSG_TRAN_INFO
	MSG_CLOSE_CONN
)

type MsgInfo struct {
	msgType uint8
	messagelen int32
	message string
}

var SyncMsgInfoPool = sync.Pool{
	New:func() interface{} {
		return &MsgInfo{}
	},
}

func (m *MsgInfo) SetMessage(msgType uint8, message string) {
	 m.msgType = msgType
	 m.messagelen = int32(len(message))
	 m.message = message
}

func (m *MsgInfo)Pack() []byte{
	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian,m.msgType)
	binary.Write(&buf, binary.LittleEndian,m.messagelen)
	buf.Write([]byte(m.message))
	SyncMsgInfoPool.Put(m)
	return buf.Bytes()
}



type MsgConnInfo struct {
	msgType uint8
	connId  int32
	messagelen int32
	message string
}

var SyncMsgConnInfoPool = sync.Pool{
	New:func() interface{} {
		return &MsgConnInfo{}
	},
}

func (m *MsgConnInfo) SetMessage(msgType uint8,connId int32, message string) {
	m.msgType = msgType
	m.connId = connId
	m.messagelen = int32(len(message))
	m.message = message
}

func (m *MsgConnInfo)Pack() ([]byte, error ){
    var buf bytes.Buffer
    err := binary.Write(&buf, binary.LittleEndian, m.msgType)
    if err != nil{
    	return nil,err
	}
	err = binary.Write(&buf, binary.LittleEndian, m.connId)
	if err != nil{
		return nil,err
	}
	//这种类型的消息不包含message
	if m.msgType == MSG_CLOSE_CONN{
		return buf.Bytes(), nil
	}
	err = binary.Write(&buf, binary.LittleEndian, m.messagelen)
	if err != nil{
		return nil,err
	}
	buf.Write([]byte(m.message))
	return buf.Bytes(), nil
}

