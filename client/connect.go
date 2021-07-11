package main

import (
	"net"
	"time"
)


func connectServer(ServerAddr string) {
	conn, err := net.DialTimeout("tcp", ServerAddr, 60*time.Second)
	if err != nil {
		return
	}

}