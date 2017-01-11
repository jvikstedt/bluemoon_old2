package main

import (
	"log"
	"net"
	"net/http"

	"github.com/gorilla/websocket"

	"github.com/jvikstedt/bluemoon/bm"
	"github.com/jvikstedt/bluemoon/socket"
	"github.com/jvikstedt/bluemoon/store"
	"github.com/jvikstedt/bluemoon/ws"
)

var hub *Hub

func main() {
	workerStore := store.NewClientStore()
	userStore := store.NewClientStore()
	userInfoStore := store.NewUserInfoStore()

	utilController := NewUtilController(userInfoStore)

	dataRouter := bm.NewDataRouter()
	dataRouter.Register("quit", utilController.Quit)
	dataRouter.Register("ping", utilController.Ping)

	hub = NewHub(dataRouter, workerStore, userStore, userInfoStore)

	sServer := socket.NewServer(manageConn)
	go sServer.Listen(":5000")

	wsServer := ws.NewServer(manageWSConn)
	http.Handle("/", wsServer)
	log.Fatal(http.ListenAndServe(":4000", nil))
}

func manageConn(conn *net.TCPConn) error {
	cw := socket.NewConnectionWrapper(conn)
	defer cw.Close()

	return hub.ManageWorkerConn(cw)
}

func manageWSConn(conn *websocket.Conn) error {
	cw := ws.NewConnectionWrapper(conn)
	defer cw.Close()

	return hub.ManageUserConn(cw)
}
