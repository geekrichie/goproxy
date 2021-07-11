package proxy

import (
	"goproxy/mux"
	"log"
	"net"
)

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func TcpProxy(address string,targetAddrs []string) {
	var err error
	listener,err := net.Listen("tcp", address)
	checkError(err)
	for{
		conn,err := listener.Accept()
		if err != nil {
           conn.Close()
		}
		connection := mux.NewConn(conn)
		connection.Target.TargetAddrs = targetAddrs
		go handleConnection(connection)
	}
}


func handleConnection(conn mux.Connection) {
	targetAddr := conn.Target.GetRandomAddr()
	conn.SendLinkInfo(targetAddr)
}

