package controller

import (
	"encoding/json"
	"fmt"

	"github.com/jvikstedt/bluemoon/bm"
	"github.com/jvikstedt/bluemoon/worker"
	"github.com/jvikstedt/bluemoon/worker/event"
)

type GateController struct {
	hub        *worker.Hub
	room       worker.Room
	userRouter *bm.DataRouter
}

func NewGateController(hub *worker.Hub, room worker.Room, userRouter *bm.DataRouter) *GateController {
	return &GateController{
		hub:        hub,
		room:       room,
		userRouter: userRouter,
	}
}

type UserEvent struct {
	UserID int `json:"user_id"`
}

func (gc *GateController) UserJoined(client bm.Client, data []byte) {
	var userEvent UserEvent
	err := json.Unmarshal(data, &userEvent)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(userEvent)
	user := worker.NewUser(userEvent.UserID)

	gc.hub.Broadcast(struct {
		Name string `json:"name"`
		ID   int    `json:"id"`
		X    int    `json:"x"`
		Y    int    `json:"y"`
	}{"new_player", user.ID(), 50, 50})
	gc.hub.AddUser(user)

	// Event based
	userJoinedEvent := &event.UserJoined{
		ID: userEvent.UserID,
	}
	gc.room.AddEvent(userJoinedEvent)
}

func (gc *GateController) UserLeft(client bm.Client, data []byte) {
	var userEvent UserEvent
	err := json.Unmarshal(data, &userEvent)
	if err != nil {
		fmt.Println(err)
	}
	gc.hub.RemoveUserByID(userEvent.UserID)

	// Event based
	userLeftEvent := &event.UserLeft{
		ID: userEvent.UserID,
	}
	gc.room.AddEvent(userLeftEvent)
}

type UserData struct {
	Name string `json:"name"`
}

func (gc *GateController) FromUser(client bm.Client, data []byte) {
	var userData UserData
	err := json.Unmarshal(data, &userData)
	if err != nil {
		fmt.Println(err)
	}

	handle, err := gc.userRouter.Handler(userData.Name)
	if err != nil {
		fmt.Println(err)
		return
	}
	handle(client, data)
}
