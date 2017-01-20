package main

import (
	"net"
	"net/http"
	"os"

	"github.com/gorilla/websocket"

	"github.com/jvikstedt/bluemoon/bm"
	"github.com/jvikstedt/bluemoon/gate"
	"github.com/jvikstedt/bluemoon/gate/controller"
	"github.com/jvikstedt/bluemoon/gate/store"
	"github.com/jvikstedt/bluemoon/logger"
	"github.com/jvikstedt/bluemoon/net/socket"
	"github.com/jvikstedt/bluemoon/net/ws"
)

var hub *gate.Hub

func main() {
	log := logger.NewLogrusLogger(os.Stdout, logger.DebugLevel)

	workerStore := store.NewClientStore()
	userStore := store.NewClientStore()
	userInfoStore := store.NewUserInfoStore()

	userController := controller.NewUserController(userInfoStore, userStore)
	workerController := controller.NewWorkerController(userInfoStore, userStore)

	userRouter := bm.NewDataRouter()
	workerRouter := bm.NewDataRouter()

	userRouter.SetDefaultHandler(userController.ToWorker)

	workerRouter.Register("to_users", workerController.ToUsers)

	hub = gate.NewHub(log, userRouter, workerRouter, workerStore, userStore, userInfoStore)

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
