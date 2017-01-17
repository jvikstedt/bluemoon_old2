package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/jvikstedt/bluemoon/bm"
	"github.com/jvikstedt/bluemoon/logger"
	"github.com/jvikstedt/bluemoon/net/socket"
	"github.com/jvikstedt/bluemoon/worker"
	"github.com/jvikstedt/bluemoon/worker/controller"
)

type DN struct {
	Name string `json:"name"`
}

func main() {
	log := logger.NewLogrusLogger(os.Stdout, logger.DebugLevel)

	sClient := socket.NewClient()
	conn, err := sClient.Connect("gate:5000")
	if err != nil {
		panic(err)
	}

	cw := socket.NewConnectionWrapper(conn)
	defer cw.Close()

	hub := worker.NewHub(nil)
	room := worker.NewRoom(log, hub)
	userController := controller.NewUserController(hub, room)

	dataRouter := bm.NewDataRouter()
	dataRouter.Register("user_joined", userController.UserJoined)
	dataRouter.Register("user_left", userController.UserLeft)
	dataRouter.Register("direction", userController.Direction)

	gate := bm.NewBaseClient(1, cw, func(client bm.Client, data []byte) {
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

	hub.SetGate(gate)

	go gate.EnableReader()
	go gate.EnableWriter()
	room.Run()
}
