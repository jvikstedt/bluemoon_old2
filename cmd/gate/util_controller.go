package main

import (
	"fmt"

	"github.com/jvikstedt/bluemoon/bm"
)

type UtilController struct {
	uis bm.UserInfoStore
}

func NewUtilController(uis bm.UserInfoStore) *UtilController {
	return &UtilController{
		uis: uis,
	}
}

func (uh *UtilController) Quit(client bm.Client, data []byte) {
	fmt.Printf("Quitting client: %d\n", client.ID())
	client.Close()
}

func (uh *UtilController) Ping(client bm.Client, data []byte) {
	userInfo := uh.uis.ByID(client.ID())
	worker := userInfo.Worker()
	worker.Write(data)

	fmt.Printf("Received ping from client: %d\n", client.ID())
	client.Write([]byte("pong"))
}
