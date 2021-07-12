package main

import (
	"flag"
	"goproxy/common"
	"goproxy/server"
	"strconv"
)

var ServerPort int

func init(){
	flag.IntVar(&ServerPort, "serverport", 0, "listening port for client to connect ")
}

func main() {
	flag.Parse()
	if ServerPort ==0 {
		ServerPort  = common.GetServerPort()
	}
	server.StartServer(":"+strconv.Itoa(ServerPort))
}