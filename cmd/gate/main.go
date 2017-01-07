package main

import (
	"fmt"
	"net"
	"net/http"

	"github.com/gorilla/websocket"

	"github.com/jvikstedt/bluemoon/gate"
	"github.com/jvikstedt/bluemoon/socket"
	"github.com/jvikstedt/bluemoon/ws"
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

func manageWSConn(conn *websocket.Conn) error {
	cw := ws.NewConnectionWrapper(conn)

	w := gate.NewUser(1, cw, func(user *gate.User, data []byte) {
		fmt.Printf("Got message from user: %s\n", string(data))
		user.Write([]byte("Pong"))
	})
	go w.EnableReader()
	w.EnableWriter()

	return nil
}

func main() {
	sServer := socket.NewServer(manageConn)
	go sServer.Listen(":5000")

	wsServer := ws.NewServer(manageWSConn)
	http.Handle("/", wsServer)
	http.ListenAndServe(":4000", nil)
}
