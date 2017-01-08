package main

import (
	"log"
	"net"
	"net/http"

	"github.com/gorilla/websocket"

	"github.com/jvikstedt/bluemoon/gate"
	"github.com/jvikstedt/bluemoon/socket"
	"github.com/jvikstedt/bluemoon/ws"
)

var hub *gate.Hub
var handler *gate.Handler

func manageConn(conn *net.TCPConn) error {
	cw := socket.NewConnectionWrapper(conn)
	defer cw.Close()

	w := gate.NewWorker(1, cw, handler.HandleWorkerData)
	hub.AddWorker(w)
	defer hub.RemoveWorker(w)

	go w.EnableReader()
	w.EnableWriter()

	return nil
}

func manageWSConn(conn *websocket.Conn) error {
	cw := ws.NewConnectionWrapper(conn)

	u := gate.NewUser(1, cw, handler.HandleUserData)
	hub.AddUser(u)
	defer hub.RemoveUser(u)

	go u.EnableReader()
	u.EnableWriter()

	return nil
}

func main() {
	hub = gate.NewHub()
	handler = gate.NewHandler(hub)

	sServer := socket.NewServer(manageConn)
	go sServer.Listen(":5000")

	wsServer := ws.NewServer(manageWSConn)
	http.Handle("/", wsServer)
	log.Fatal(http.ListenAndServe(":4000", nil))
}
