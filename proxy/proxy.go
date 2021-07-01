package proxy

import (
	"io"
	"log"
	"net"
	"time"
)

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func TcpProxy(address string ) {
	var err error
	listener,err := net.Listen("tcp", address)
	checkError(err)
	for{
		conn,err := listener.Accept()
		checkError(err)
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	targetConn , err := net.DialTimeout("tcp","10.220.162.12:13443", 30*time.Second)
	checkError(err)
	go func() {
		//fmt.Println("start copy targetConn to conn")
		io.Copy(conn, targetConn)
		//fmt.Println("end copy targetConn to conn")
	}()
	go func() {
		//fmt.Println("start copy conn to targetConn")
		io.Copy(targetConn, conn)
		//fmt.Println("end copy conn to targetConn")
	}()
}