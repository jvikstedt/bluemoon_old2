package main

import (
	"fmt"

	"github.com/jvikstedt/bluemoon/bluemoon"
)

type UtilController struct {
}

func NewUtilController() *UtilController {
	return &UtilController{}
}

func (uh *UtilController) Quit(client bluemoon.Client, data []byte) {
	fmt.Printf("Quitting client: %d\n", client.ID())
	client.Close()
}

func (uh *UtilController) Ping(client bluemoon.Client, data []byte) {
	fmt.Printf("Received ping from client: %d\n", client.ID())
	client.Write([]byte("pong"))
}
