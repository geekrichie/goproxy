package server

import (
	"fmt"
	"goproxy/file"
	"goproxy/log"
	"goproxy/mux/mux_link"
	"goproxy/mux/mux_msg"
	"goproxy/mux/mux_net"
	"net"
	"strconv"
)

type Server struct {
	cm mux_net.ConnManager
	address string
}
var defaultServer Server


func StartServer(address string) {
	var err error
	defaultServer.address = address
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
		defaultServer.cm.AddConnection(connection)
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
	log.Infof("client: %s mode : %d", conn.GetConn().RemoteAddr().String(), modebyte)
	if err != nil {
		conn.Close()
		return
	}
	mode := modebyte[0]
	//设置连接类型
	conn.SetConnType(int(mode))
	switch mode{
		case mux_link.MainMode:
			//tcpconn := conn.GetConn().(*net.TCPConn)
			//tcpconn.SetKeepAlive(true)
			//tcpconn.SetKeepAlivePeriod(5*time.Second)
			conn.SendMsg(mux_msg.MSG_LOG_INFO, "this is a first logging message")
			conn.SendMsg(mux_msg.MSG_LOG_INFO, "this is second message")
		case mux_link.TranMode:
			transConn(conn)


	}
	//targetAddr := conn.Target.GetRandomAddr()
	//conn.SendLinkInfo(targetAddr)
}

func transConn(conn mux_net.Connection) {

	taskDb := file.LoadTask()
	for _, task := range taskDb.Tasks {
		conn.Target.TargetAddrs = task.TargetAddrs
		go listenOuterConn(task)
		targetAddr := conn.Target.GetRandomAddr()
		conn.SendLinkInfo(targetAddr)
		break
	}

}

func listenOuterConn (task file.Task) {
	l, err := net.Listen("tcp", ":"+strconv.Itoa(task.Port))
	log.Info(fmt.Sprintf("New Task Listen port : %d", task.Port))
	if err != nil {
		log.Error(err.Error())
		return
	}
	for {
		netconn, err  := l.Accept()
		if  err != nil {
			log.Error(err.Error())
			netconn.Close()
		}

	}
}