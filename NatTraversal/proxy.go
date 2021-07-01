package NatTraversal

import (
	"io"
	"log"
	"net"
	"net/http"
	"sync"
	"time"
)

func Serve() {
     var mux  = http.NewServeMux()
     mux.HandleFunc("/path", func(w http.ResponseWriter, r *http.Request) {
     	io.WriteString(w,"hello world")
	 })
     fs := http.FileServer(http.Dir("../static/"))
     mux.HandleFunc("/image/", http.StripPrefix("/image/",fs).ServeHTTP)
     http.ListenAndServe(":5000", mux)
}

func LocalOperation() {
     conn, err := net.DialTimeout("tcp","10.220.169.69:10002", 30*time.Second)
     if err != nil {
     	log.Fatal(err)
	 }
     targetconn, err := net.DialTimeout("tcp", ":5000", 30*time.Second)
	if err != nil {
		log.Fatal(err)
	}
	var wg  sync.WaitGroup
     wg.Add(2)
	go func() {
		defer wg.Done()
		for {
			_, err := io.Copy(conn, targetconn)
			if err != nil {
				log.Fatal(err)
			}
		}
	}()
     go func() {
     	defer wg.Done()
     	for {
			_, err = io.Copy(targetconn, conn)
			if err != nil {
				log.Fatal(err)
			}
		}
	 }()
     wg.Wait()
}


var streams [2]net.Conn

func RemoteOperation() {
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		Listen(":10003",0)
		wg.Done()
	}()
	go func() {
		Listen(":10002", 1)
		wg.Done()
	}()
	wg.Wait()
}

func checkError( err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func Listen(addr string, seq int) {
	listener, err := net.Listen("tcp", addr)
	checkError(err)
	for{
		conn, err := listener.Accept()
		streams[seq] = conn
		checkError(err)
		go handleConnection(conn, seq)
	}
}


func handleConnection(conn net.Conn, seq int) {
	if streams[1-seq] ==  nil {
		return
	}
	go func() {
		for {
			_, err := io.Copy(conn, streams[1-seq])
			log.Println(err)
		}
	}()
	go func() {
		for {
			_, err := io.Copy(streams[1-seq], conn)
			log.Println(err)
		}
	}()
}