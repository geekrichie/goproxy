package main

import "net/http"

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/help", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world"))
	})
	fd := http.FileServer(http.Dir("./static"))
	mux.Handle("/img/",http.StripPrefix("/img/", fd))
	http.ListenAndServe(":8080", mux)
}
