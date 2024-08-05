package main

import (
	"net/http"
	"os"
)

const (
	SERVER_DEFAULT_ADDR = ":8080"
)

func handler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./static.html")
}

func serveHtml() {
	fs := http.FileServer(http.Dir("public"))

	http.Handle("/", fs)
	server_addr := os.Getenv("SERVER_ADDR")
	if server_addr == "" {
		println("Use SERVER_DEFAULT_ADDR")
		server_addr = SERVER_DEFAULT_ADDR
	}
	err := http.ListenAndServe(server_addr, nil)
	if err != nil {
		panic(err)
	}
}
