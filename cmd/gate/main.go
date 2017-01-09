package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/gorilla/websocket"

	"github.com/jvikstedt/bluemoon/bluemoon"
	"github.com/jvikstedt/bluemoon/gate/client"
	"github.com/jvikstedt/bluemoon/gate/controller"
	"github.com/jvikstedt/bluemoon/gate/store"
	"github.com/jvikstedt/bluemoon/socket"
	"github.com/jvikstedt/bluemoon/ws"
)

var userStore *store.UserStore
var workerStore *store.WorkerStore
var userRouter *bluemoon.DataRouter
var workerRouter *bluemoon.DataRouter

type DN struct {
	Name string `json:"name"`
}

func manageConn(conn *net.TCPConn) error {
	cw := socket.NewConnectionWrapper(conn)
	defer cw.Close()

	w := client.NewWorkerClient(1, cw, func(client bluemoon.Client, data []byte) {
		fmt.Printf("New message from worker: %d\n", client.ID())
		fmt.Print(string(data))
		var dn DN
		err := json.Unmarshal(data, &dn)
		if err != nil {
			fmt.Println(err)
			return
		}
		handle, err := workerRouter.Handler(dn.Name)
		if err != nil {
			fmt.Println(err)
			return
		}
		handle(client, data)
	})
	workerStore.Add(w)
	defer workerStore.Remove(w)

	go w.EnableReader()
	w.EnableWriter()

	return nil
}

func manageWSConn(conn *websocket.Conn) error {
	cw := ws.NewConnectionWrapper(conn)

	u := client.NewUserClient(1, cw, func(client bluemoon.Client, data []byte) {
		fmt.Printf("New message from user: %d\n", client.ID())
		fmt.Print(string(data))
		var dn DN
		err := json.Unmarshal(data, &dn)
		if err != nil {
			fmt.Println(err)
			return
		}
		handle, err := userRouter.Handler(dn.Name)
		if err != nil {
			fmt.Println(err)
			return
		}
		handle(client, data)
	})

	userStore.Add(u)
	defer userStore.Remove(u)

	go u.EnableReader()
	u.EnableWriter()

	return nil
}

func main() {
	userStore = store.NewUserStore()
	workerStore = store.NewWorkerStore()

	utilController := controller.NewUtilController()

	workerRouter = bluemoon.NewDataRouter()
	workerRouter.Register("quit", utilController.Quit)
	workerRouter.Register("ping", utilController.Ping)

	userRouter = bluemoon.NewDataRouter()
	userRouter.Register("quit", utilController.Quit)
	userRouter.Register("ping", utilController.Ping)

	sServer := socket.NewServer(manageConn)
	go sServer.Listen(":5000")

	wsServer := ws.NewServer(manageWSConn)
	http.Handle("/", wsServer)
	log.Fatal(http.ListenAndServe(":4000", nil))
}
