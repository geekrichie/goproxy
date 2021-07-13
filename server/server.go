package server

import (
	"goproxy/log"
	"goproxy/mux/mux_link"
	"goproxy/mux/mux_net"
	"net"
	"time"
)

type Server struct {
	cm mux_net.ConnManager
	address string
}
var server Server


func StartServer(address string) {
	var err error
	server.address = address
	listener,err := net.Listen("tcp", address)
	if err != nil {
		panic(err)
	}
	for{
		conn,err := listener.Accept()
		if err != nil {
           conn.Close()
		}
		log.Infof("client : %s connect to the server\n",conn.RemoteAddr().String())
		connection := mux_net.NewConn(conn)
		go handleConnection(connection)
	}
}


func handleConnection(conn mux_net.Connection) {
	err := conn.ReceiveHandShake()
	if err != nil {
		log.Error(err.Error())
	}else{
		log.Info("handshake finished!")
	}
	modebyte, err := conn.ReadContent(1)
	if err != nil {
		conn.Close()
		return
	}
	mode := modebyte[0]
	switch mode{
		case mux_link.MainMode:
			tcpconn := conn.GetConn().(*net.TCPConn)
			tcpconn.SetKeepAlive(true)
			tcpconn.SetKeepAlivePeriod(5*time.Second)


	}
	//targetAddr := conn.Target.GetRandomAddr()
	//conn.SendLinkInfo(targetAddr)
}

