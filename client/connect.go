package client

import (
	"goproxy/log"
	"goproxy/mux/mux_link"
	"goproxy/mux/mux_msg"
	"net"
	"goproxy/mux/mux_net"
	"time"
)




func ConnectServer(ServerAddr string) {
	startServer(ServerAddr, mux_link.MainMode)
}

func startServer(addr string, mode int) {
	conn, err := net.DialTimeout("tcp", addr, 60*time.Second)
	if err != nil {
		log.Error(err.Error())
		return
	}
	newConn := mux_net.NewConn(conn)
	err = newConn.SendHandShake()
	if err != nil {
		log.Error(err.Error())
		conn.Close()
		return
	}
	switch mode{
	case mux_link.MainMode:
		go handleMainConnect(newConn)
	case mux_link.TranMode:
		go handleTranConnect(newConn)

	}

}

func handleTranConnect(conn mux_net.Connection) {
	
}

func handleMainConnect(conn mux_net.Connection) {
	conn.SendMode(mux_link.MainMode)
	for {
		msgType ,err := conn.ReadMsgType()
		if err != nil {
			conn.Close()
			return
		}
		switch msgType{
			case mux_msg.MSG_INFO:
			    handleInfoMsg(conn)
		}

	}
}

func handleInfoMsg(conn mux_net.Connection) {
	content, err := conn.ReadLenContent()
	if err != nil {
		return
	}
	log.Info(string(content))
}
