package mux_link

import (
	"fmt"
	"testing"
	"time"
)


func TestReceiveWindow_Read(t *testing.T) {
	rw := NewReceiveWindow()
	var b = make([]byte, 10)
	go func() {
		n, err := rw.Write([]byte("abcde"))
		if err != nil {
			t.Error(err)
		}
		t.Logf("write %d byte\n", n)
		n, err = rw.Write([]byte("abcdef"))
		if err != nil {
			t.Error(err)
		}
		t.Logf("write %d byte\n", n)
		time.Sleep(2*time.Second)
		rw.Write([]byte("kef"))
	}()
	time.Sleep(1*time.Second)
	n, err := rw.Read(b)
	if err != nil {
		t.Error(err)
	}
	t.Log(n)
	fmt.Printf("%#v\n", b)
	fmt.Printf("%#v\n", rw)
	c := make([]byte,2)
	n, err = rw.Read(c)
	if err != nil {
		t.Error(err)
	}
	t.Log(n)
	fmt.Printf("%#v\n", c)
	n, err = rw.Read(c)
	if err != nil {
		t.Error(err)
	}
	t.Log(n)
	fmt.Printf("%#v\n", c)
}

func TestSlice(t *testing.T) {

}