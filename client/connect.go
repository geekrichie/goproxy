package client

import (
	log "github.com/amoghe/distillog"
	"goproxy/mux"
	"net"
	"time"
)


func ConnectServer(ServerAddr string) {
	conn, err := net.DialTimeout("tcp", ServerAddr, 60*time.Second)
	if err != nil {
		log.Errorln(err)
		return
	}
	defer conn.Close()
	newConn := mux.NewConn(conn)
	err = newConn.SendHandShake()
	if err != nil {
		log.Errorln(err)
		return
	}
}