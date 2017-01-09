package controller

import (
	"fmt"

	"github.com/jvikstedt/bluemoon/gate/client"
)

type UtilController struct {
}

func NewUtilController() *UtilController {
	return &UtilController{}
}

func (uh *UtilController) Quit(user *client.User, data []byte) {
	fmt.Printf("Quitting client: %d\n", user.ID())
	user.Close()
}

func (uh *UtilController) Ping(user *client.User, data []byte) {
	fmt.Printf("Received ping from client: %d\n", user.ID())
	user.Write([]byte("pong"))
}
