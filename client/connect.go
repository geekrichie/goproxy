package client

import (
	"goproxy/mux"
	"net"
	"time"
)


func ConnectServer(ServerAddr string) {
	conn, err := net.DialTimeout("tcp", ServerAddr, 60*time.Second)
	if err != nil {
		return
	}
	newConn := mux.NewConn(conn)
	newConn.SendHandShake()
}