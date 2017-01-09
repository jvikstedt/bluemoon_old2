package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/gorilla/websocket"

	"github.com/jvikstedt/bluemoon/bluemoon"
	"github.com/jvikstedt/bluemoon/gate"
	"github.com/jvikstedt/bluemoon/gate/controller"
	"github.com/jvikstedt/bluemoon/socket"
	"github.com/jvikstedt/bluemoon/ws"
)

var hub *gate.Hub
var dataRouter *bluemoon.DataRouter

type DN struct {
	Name string `json:"name"`
}

func manageConn(conn *net.TCPConn) error {
	cw := socket.NewConnectionWrapper(conn)
	defer cw.Close()

	w := bluemoon.NewBaseClient(1, cw, func(client bluemoon.Client, data []byte) {
		fmt.Printf("New message from worker: %d\n", client.ID())
		fmt.Print(string(data))
		var dn DN
		err := json.Unmarshal(data, &dn)
		if err != nil {
			fmt.Println(err)
			return
		}
		handle, err := dataRouter.Handler(dn.Name)
		if err != nil {
			fmt.Println(err)
			return
		}
		handle(client, data)
	})
	hub.AddWorker(w)
	defer hub.RemoveWorker(w)

	go w.EnableReader()
	w.EnableWriter()

	return nil
}

func manageWSConn(conn *websocket.Conn) error {
	cw := ws.NewConnectionWrapper(conn)

	u := bluemoon.NewBaseClient(1, cw, func(client bluemoon.Client, data []byte) {
		fmt.Printf("New message from user: %d\n", client.ID())
		fmt.Print(string(data))
		var dn DN
		err := json.Unmarshal(data, &dn)
		if err != nil {
			fmt.Println(err)
			return
		}
		handle, err := dataRouter.Handler(dn.Name)
		if err != nil {
			fmt.Println(err)
			return
		}
		handle(client, data)
	})

	hub.AddUser(u)
	defer hub.RemoveUser(u)

	go u.EnableReader()
	u.EnableWriter()

	return nil
}

func main() {
	utilController := controller.NewUtilController()

	dataRouter = bluemoon.NewDataRouter()
	dataRouter.Register("quit", utilController.Quit)
	dataRouter.Register("ping", utilController.Ping)

	hub = gate.NewHub()

	sServer := socket.NewServer(manageConn)
	go sServer.Listen(":5000")

	wsServer := ws.NewServer(manageWSConn)
	http.Handle("/", wsServer)
	log.Fatal(http.ListenAndServe(":4000", nil))
}
