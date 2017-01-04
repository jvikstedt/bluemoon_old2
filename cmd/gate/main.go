package main

import (
	"net/http"
)

func main() {
	go listenSockets(":5000")

	http.HandleFunc("/", HandleWS)
	http.ListenAndServe(":4000", nil)
}
