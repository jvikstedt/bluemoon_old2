package main

import (
	"fmt"

	"github.com/jvikstedt/bluemoon/bluemoon"
)

type UtilController struct {
	uis bluemoon.UserInfoStore
}

func NewUtilController(uis bluemoon.UserInfoStore) *UtilController {
	return &UtilController{
		uis: uis,
	}
}

func (uh *UtilController) Quit(client bluemoon.Client, data []byte) {
	fmt.Printf("Quitting client: %d\n", client.ID())
	client.Close()
}

func (uh *UtilController) Ping(client bluemoon.Client, data []byte) {
	userInfo := uh.uis.ByID(client.ID())
	worker := userInfo.Worker()
	worker.Write(data)

	fmt.Printf("Received ping from client: %d\n", client.ID())
	client.Write([]byte("pong"))
}
