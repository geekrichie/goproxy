package main

import (
	"encoding/binary"
	"fmt"
	"testing"
)

type structa struct {
	b *structb
}
type structb struct {
	 Num int
}

func (b *structb ) Add () {
	b.Num = b.Num +1
}

func Test_Net(t *testing.T) {
	//conn, err := net.Dial("tcp",":9998")
	//if err != nil {
	//	t.Log(err)
	//}
	//conn.Write([]byte("hello"))
	var c = structa{}
	var d = structb{}
	c.b = &d
	var a = c
	c.b.Add()
	c.b.Add()
	fmt.Println(a.b.Num)
}


func TestBinary(t *testing.T) {
	var buf  = []byte{uint8(14),uint8(0),uint8(0),uint8(0)}
	fmt.Println(buf)
	var message uint32
	message = binary.LittleEndian.Uint32(buf)
	fmt.Printf("bigendian enc = %x\n", message)
}