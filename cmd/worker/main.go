package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/jvikstedt/bluemoon/bm"
	"github.com/jvikstedt/bluemoon/logger"
	"github.com/jvikstedt/bluemoon/net/socket"
	"github.com/jvikstedt/bluemoon/worker"
	"github.com/jvikstedt/bluemoon/worker/action"
	"github.com/jvikstedt/bluemoon/worker/controller"
	"github.com/jvikstedt/bluemoon/worker/room"
)

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
	room := room.NewGame(log, 1, hub)
	room.AddAction(&action.AppleAction{})

	userController := controller.NewUserController(hub, room)

	userRouter := bm.NewDataRouter()
	userRouter.Register("change_dir", userController.ChangeDir)

	gateController := controller.NewGateController(hub, room, userRouter)
	gateRouter := bm.NewDataRouter()
	gateRouter.Register("user_joined", gateController.UserJoined)
	gateRouter.Register("user_left", gateController.UserLeft)
	gateRouter.Register("from_user", gateController.FromUser)

	gate := bm.NewBaseClient(1, cw, func(client bm.Client, data []byte) {
		fmt.Printf("New message from gate: %d\n", client.ID())
		fmt.Print(string(data))
		var gateIn worker.GateIn
		err := json.Unmarshal(data, &gateIn)
		if err != nil {
			fmt.Println(err)
			return
		}
		handle, err := gateRouter.Handler(gateIn.Name)
		if err != nil {
			fmt.Println(err)
			return
		}
		handle(client, *gateIn.Payload)
	})

	hub.SetGate(gate)

	go gate.EnableReader()
	go gate.EnableWriter()
	room.Run()
}
