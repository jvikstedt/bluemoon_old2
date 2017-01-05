package main

import (
	"fmt"
	"net"
	"net/http"

	"github.com/jvikstedt/bluemoon/gate"
	"github.com/jvikstedt/bluemoon/socket"
)

func manageConn(conn *net.TCPConn) error {
	cw := socket.NewConnectionWrapper(conn)
	defer cw.Close()

	w := gate.NewWorker(1, cw, func(worker *gate.Worker, data []byte) {
		fmt.Printf("Got message: %s\n", string(data))
		worker.Write([]byte("Pong"))
	})
	go w.EnableReader()
	w.EnableWriter()

	return nil
}

func main() {
	sServer := socket.NewServer(manageConn)
	go sServer.Listen(":5000")

	http.HandleFunc("/", HandleWS)
	http.ListenAndServe(":4000", nil)
}
