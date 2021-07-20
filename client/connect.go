package client

import (
	"fmt"
	"goproxy/log"
	"goproxy/mux/mux_link"
	"goproxy/mux/mux_msg"
	"goproxy/mux/mux_net"
	"net"
	"time"
)




func ConnectServer(ServerAddr string) {
	startConnect(ServerAddr, mux_link.MainMode)
}

func startConnect(addr string, mode uint8) {
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
		 handleMainConnect(newConn)
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
			case mux_msg.MSG_LOG_INFO:
			    handleInfoMsg(conn)
		}

	}
}

func handleInfoMsg(conn mux_net.Connection) {
	content, err := conn.ReadLenContent()
	if err != nil {
		fmt.Println(err)
		return
	}
	log.Info(string(content))
}
