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
)

type MsgInfo struct {
	msgType uint8
	messagelen int
	message string
}

var syncMsgInfoPool = sync.Pool{
	New:func() interface{} {
		return MsgInfo{}
	},
}

func pack(msgType uint8, message string) []byte{
	msgInfo := syncMsgInfoPool.Get().(MsgInfo)
	msgInfo.msgType = msgType
	msgInfo.messagelen = len(message)
	msgInfo.message = message
	var buf bytes.Buffer
	binary.Write(&buf,binary.LittleEndian, &msgInfo)
	syncMsgInfoPool.Put(msgInfo)
	return buf.Bytes()
}

func unpack(msg []byte) MsgInfo{
	msgInfo := *(*MsgInfo)(unsafe.Pointer(&msg))
	return msgInfo
}