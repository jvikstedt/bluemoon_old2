package main

import (
	"net/http"

	"github.com/jvikstedt/bluemoon/socket"
)

func main() {
	sServer := socket.NewServer()
	go sServer.Listen(":5000")

	http.HandleFunc("/", HandleWS)
	http.ListenAndServe(":4000", nil)
}
