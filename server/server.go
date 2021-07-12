package server

import (
	"goproxy/log"
	"goproxy/mux"
	"net"
)


func StartServer(address string) {
	var err error
	listener,err := net.Listen("tcp", address)
	if err != nil {
		panic(err)
	}
	for{
		conn,err := listener.Accept()
		if err != nil {
           conn.Close()
		}
		log.Infof(" client : %s connect to the server\n",conn.RemoteAddr().String())
		connection := mux.NewConn(conn)
		go handleConnection(connection)
	}
}


func handleConnection(conn mux.Connection) {
	conn.ReceiveHandShake()
	//targetAddr := conn.Target.GetRandomAddr()
	//conn.SendLinkInfo(targetAddr)
}

