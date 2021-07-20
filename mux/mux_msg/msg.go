package mux_msg

import (
	"bytes"
	"encoding/binary"
	"sync"
	"unsafe"
)

const (
	MSG_PING uint8 = iota
	MSG_LOG_INFO
	MSG_LINK_INFO
)

type MsgInfo struct {
	msgType uint8
	messagelen int32
	message string
}

var syncMsgInfoPool = sync.Pool{
	New:func() interface{} {
		return MsgInfo{}
	},
}

func Pack(msgType uint8, message string) []byte{
	msgInfo := syncMsgInfoPool.Get().(MsgInfo)
	msgInfo.msgType = msgType
	msgInfo.messagelen = int32(len(message))
	msgInfo.message = message
	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian,msgInfo.msgType)
	binary.Write(&buf, binary.LittleEndian,msgInfo.messagelen)
	buf.Write([]byte(msgInfo.message))
	syncMsgInfoPool.Put(msgInfo)
	return buf.Bytes()
}


func Unpack(msg []byte) MsgInfo{
	msgInfo := *(*MsgInfo)(unsafe.Pointer(&msg))
	return msgInfo
}