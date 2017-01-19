package main

import (
	"net"
	"net/http"

	"github.com/gorilla/websocket"

	"github.com/jvikstedt/bluemoon/bm"
	"github.com/jvikstedt/bluemoon/gate"
	"github.com/jvikstedt/bluemoon/gate/controller"
	"github.com/jvikstedt/bluemoon/gate/store"
	"github.com/jvikstedt/bluemoon/net/socket"
	"github.com/jvikstedt/bluemoon/net/ws"
)

var hub *gate.Hub

func main() {
	workerStore := store.NewClientStore()
	userStore := store.NewClientStore()
	userInfoStore := store.NewUserInfoStore()

	userController := controller.NewUserController(userInfoStore, userStore)

	userRouter := bm.NewDataRouter()
	workerRouter := bm.NewDataRouter()

	userRouter.SetDefaultHandler(userController.ToWorker)

	hub = gate.NewHub(userRouter, workerRouter, workerStore, userStore, userInfoStore)

	sServer := socket.NewServer(manageConn)
	go sServer.Listen(":5000")

	wsServer := ws.NewServer(manageWSConn)
	http.Handle("/", wsServer)
	http.ListenAndServe(":4000", nil)
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
