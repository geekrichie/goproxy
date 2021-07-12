package main

import (
	"flag"
	"goproxy/client"
)

var ServerAddr string

func init() {
	flag.StringVar(&ServerAddr, "serveraddr", "","default connect server addr")
}

func main() {
	//启动两个协程，主协程用于获取
	flag.Parse()
	run()
}

func run() {
	if ServerAddr == "" {
         panic("serveraddr cannot be null")
	}
	client.ConnectServer(ServerAddr)
}