package main

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/gorilla/websocket"

	"github.com/jvikstedt/bluemoon/gate"
	"github.com/jvikstedt/bluemoon/socket"
	"github.com/jvikstedt/bluemoon/ws"
)

var hub *gate.Hub

func manageConn(conn *net.TCPConn) error {
	cw := socket.NewConnectionWrapper(conn)
	defer cw.Close()

	w := gate.NewWorker(1, cw, func(worker *gate.Worker, data []byte) {
		fmt.Printf("Got message: %s\n", string(data))
		worker.Write([]byte("Pong"))
	})
	hub.AddWorker(w)
	defer hub.RemoveWorker(w)

	go w.EnableReader()
	w.EnableWriter()

	return nil
}

func manageWSConn(conn *websocket.Conn) error {
	cw := ws.NewConnectionWrapper(conn)

	u := gate.NewUser(1, cw, func(user *gate.User, data []byte) {
		fmt.Printf("Got message from user: %s\n", string(data))
		user.Write([]byte("Pong"))
	})
	hub.AddUser(u)
	defer hub.RemoveUser(u)

	go u.EnableReader()
	u.EnableWriter()

	return nil
}

func main() {
	hub = gate.NewHub()

	sServer := socket.NewServer(manageConn)
	go sServer.Listen(":5000")

	wsServer := ws.NewServer(manageWSConn)
	http.Handle("/", wsServer)
	log.Fatal(http.ListenAndServe(":4000", nil))
}
