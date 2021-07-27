package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"testing"
)

func Test_Net(t *testing.T) {
	conn, err := net.Dial("tcp",":9998")
	if err != nil {
		t.Log(err)
	}
	conn.Write([]byte("hello"))
}


func TestBinary(t *testing.T) {
	var buf  = []byte{uint8(14),uint8(0),uint8(0),uint8(0)}
	fmt.Println(buf)
	var message uint32
	message = binary.LittleEndian.Uint32(buf)
	fmt.Printf("bigendian enc = %x\n", message)
}