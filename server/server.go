package server

import (
	"encoding/binary"
	"fmt"
	"goproxy/file"
	"goproxy/log"
	"goproxy/mux/mux_link"
	"goproxy/mux/mux_msg"
	"goproxy/mux/mux_net"
	"io"
	"net"
	"strconv"
)

type Server struct {
	cm mux_net.ConnManager
	address string
	plexer *mux_link.MultiPlexer
	taskDb   *file.Db
}
var defaultServer Server


func StartServer(address string) {
	var err error
	defaultServer.address = address
	LoadTaskJob()
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
		conn.Close()
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
			defaultServer.plexer = mux_link.NewMultiPlexer(conn.GetConn())
			conn.Plexer = defaultServer.plexer
			go transConn(conn)


	}
	//targetAddr := conn.Target.GetRandomAddr()
	//conn.SendLinkInfo(targetAddr)
}

func transConn(conn mux_net.Connection) {

	for {
		msgType ,err := conn.ReadMsgType()
		if err != nil {
			conn.Close()
			return
		}
		switch msgType{
			case mux_msg.MSG_TRAN_INFO:
				log.Info("accept msg_tran_info")
				Unpack(conn)

		}
	}

}


func Unpack(conn mux_net.Connection) {
	var buf = make([]byte, 4)
	io.ReadFull(&conn, buf)
	var connId uint32
	connId = binary.LittleEndian.Uint32(buf)
	linkConn := conn.Plexer.GetConnById(int(connId))
	io.ReadFull(&conn, buf)
	var messagelen uint32
	messagelen = binary.LittleEndian.Uint32(buf)
	log.Infof("New messagelen : %d", messagelen)
	var message = make([]byte, messagelen)
	io.ReadFull(&conn, message)
	linkConn.ReceiveWindowWrite(message)
	return
}

func LoadTaskJob() {
	taskDb := file.LoadTask()
	defaultServer.taskDb = taskDb
	for _, task := range taskDb.Tasks {
		go listenOuterConn(task, ProxyTcpConnect)
		break
	}
}
func ProxyTcpConnect(proxyconn net.Conn, task file.Task) {
	plexerConn := mux_link.NewConn(defaultServer.plexer)
	defaultServer.plexer.AddConn(plexerConn)
	target := mux_net.Target{TargetAddrs: task.TargetAddrs}
	plexerConn.SendLinkInfo(target.GetRandomAddr())
	mux_link.Copy(proxyconn, plexerConn)
}

func listenOuterConn (task file.Task,f func(net.Conn, file.Task) ) {
	l, err := net.Listen("tcp", ":"+strconv.Itoa(task.Port))
	log.Info(fmt.Sprintf("New Task Listen port : %d", task.Port))
	if err != nil {
		log.Error(err.Error())
		return
	}
	for {
		netconn, err  := l.Accept()
		log.Infof("client %s connect to the port %d", netconn.RemoteAddr().String(), task.Port)
		if  err != nil {
			log.Error(err.Error())
			netconn.Close()
		}
		go f(netconn,task)
	}
}