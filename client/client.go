package main

import (
	"errors"
	"flag"
)

var ServerAddr string

func init() {
   flag.StringVar(&ServerAddr, "serveraddr", "", "connect to server's address")
}

func main() {
	flag.Parse()
	runClient()
}

func runClient() error{
	 if ServerAddr == "" {
	 	return errors.New("server addr cannot be null")
	 }
	 connectServer(ServerAddr)
	 return nil
}