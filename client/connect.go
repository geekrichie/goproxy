package client

import (
	"encoding/binary"
	"goproxy/log"
	"goproxy/mux/mux_link"
	"goproxy/mux/mux_msg"
	"goproxy/mux/mux_net"
	"io"
	"net"
	"time"
)




func ConnectServer(ServerAddr string) {
	go startConnect(ServerAddr, mux_link.TranMode)
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
	conn.SendMode(mux_link.TranMode)
	conn.Plexer  = mux_link.NewMultiPlexer(conn.GetConn())
	for {
		msgType ,err := conn.ReadMsgType()
		if err != nil {
			conn.Close()
			return
		}
		switch msgType{
		case mux_msg.MSG_LINK_INFO:
			//log.Info("accept msg_link_info")
			dealNewTaskConn(conn)
		case mux_msg.MSG_TRAN_INFO:
			//log.Info("accept msg_tran_info")
			mux_net.Unpack(conn, mux_msg.MSG_TRAN_INFO)
		}
	}
}




func dealNewTaskConn(conn mux_net.Connection) {
	linkConn := mux_link.NewConn(conn.Plexer)
	var buf = make([]byte, 4)
	io.ReadFull(&conn, buf)
	var connId uint32
	connId = binary.LittleEndian.Uint32(buf)
	linkConn.SetConnId(int(connId))
	//log.Infof("New conn Id : %d", connId)
	conn.Plexer.AddConn(linkConn)

	io.ReadFull(&conn, buf)
	var messagelen uint32
	messagelen = binary.LittleEndian.Uint32(buf)
	//log.Infof("New messagelen : %d", messagelen)
	var linkinfo = make([]byte, messagelen)
	io.ReadFull(&conn, linkinfo)

	c, err := net.DialTimeout("tcp", string(linkinfo), time.Millisecond*200)
	if err != nil {
		log.Errorf("dial %s error", string(linkinfo))
	}

	mux_link.Copy(c, linkConn)

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
		log.Error(err.Error())
		return
	}
	log.Info(string(content))
}
