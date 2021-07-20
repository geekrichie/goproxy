package mux_msg

import (
	"fmt"
	"testing"
)

func TestPack(t *testing.T) {
	fmt.Println('t')
	buf := Pack(MSG_LOG_INFO, "this is the first message")
	fmt.Println(buf)
}