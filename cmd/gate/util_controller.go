package main

import (
	"encoding/json"
	"fmt"

	"github.com/jvikstedt/bluemoon/bm"
)

type UtilController struct {
	uis       bm.UserInfoStore
	userStore bm.ClientStore
}

func NewUtilController(uis bm.UserInfoStore, userStore bm.ClientStore) *UtilController {
	return &UtilController{
		uis:       uis,
		userStore: userStore,
	}
}

func (uh *UtilController) Quit(client bm.Client, data []byte) {
	fmt.Printf("Quitting client: %d\n", client.ID())
	client.Close()
}

func (uh *UtilController) Ping(client bm.Client, data []byte) {
	userInfo, err := uh.uis.ByID(client.ID())
	if err != nil {
		return
	}
	worker := userInfo.Worker()
	worker.Write(data)

	fmt.Printf("Received ping from client: %d\n", client.ID())
	client.Write([]byte("pong"))
}

type UserData struct {
	UserID int `json:"user_id"`
}

func (uh *UtilController) Move(client bm.Client, data []byte) {
	var userData UserData
	err := json.Unmarshal(data, &userData)
	if err != nil {
		fmt.Println(err)
		return
	}
	user, err := uh.userStore.ByID(userData.UserID)
	if err != nil {
		fmt.Println(err)
		return
	}
	user.Write(data)
}

func (uh *UtilController) Direction(client bm.Client, data []byte) {
	userInfo, err := uh.uis.ByID(client.ID())
	if err != nil {
		fmt.Println(err)
		return
	}
	worker := userInfo.Worker()
	worker.Write([]byte(fmt.Sprintf(`{"name": "direction", "user_id": %d, "payload": %s}`, client.ID(), string(data)) + "\n"))
}
