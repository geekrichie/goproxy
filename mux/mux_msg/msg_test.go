package mux_msg

import (
	"fmt"
	"log"
	"testing"
)

func TestPack(t *testing.T) {
	msgInfo := SyncMsgInfoPool.Get().(*MsgInfo)
	msgInfo.SetMessage(MSG_LINK_INFO, "this is a new message")
	buf := msgInfo.Pack()
	fmt.Println(buf)
}

func TestMsgConnInfo_Pack(t *testing.T) {
	msgInfo := SyncMsgConnInfoPool.Get().(*MsgConnInfo)
	msgInfo.SetMessage(MSG_LINK_INFO,12, "this is a new message")
	buf,err := msgInfo.Pack()
	if err != nil {
		log.Println(err)
	}
	fmt.Println(buf)
}