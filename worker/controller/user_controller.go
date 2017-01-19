package controller

import (
	"encoding/json"
	"fmt"

	"github.com/jvikstedt/bluemoon/bm"
	"github.com/jvikstedt/bluemoon/worker"
	"github.com/jvikstedt/bluemoon/worker/event"
)

type UserController struct {
	hub  *worker.Hub
	room worker.Room
}

func NewUserController(hub *worker.Hub, room worker.Room) *UserController {
	return &UserController{
		hub:  hub,
		room: room,
	}
}

type UserEvent struct {
	Name    string `json:"name"`
	Payload struct {
		UserID int `json:"user_id"`
	} `json:"payload"`
}

func (uc *UserController) UserJoined(client bm.Client, data []byte) {
	var userEvent UserEvent
	err := json.Unmarshal(data, &userEvent)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(userEvent)
	user := worker.NewUser(userEvent.Payload.UserID)

	uc.hub.Broadcast(struct {
		Name string `json:"name"`
		ID   int    `json:"id"`
		X    int    `json:"x"`
		Y    int    `json:"y"`
	}{"new_player", user.ID(), 50, 50})
	uc.hub.AddUser(user)

	// Event based
	userJoinedEvent := &event.UserJoined{
		ID: userEvent.Payload.UserID,
	}
	uc.room.AddEvent(userJoinedEvent)
}

func (uc *UserController) UserLeft(client bm.Client, data []byte) {
	var userEvent UserEvent
	err := json.Unmarshal(data, &userEvent)
	if err != nil {
		fmt.Println(err)
	}
	uc.hub.RemoveUserByID(userEvent.Payload.UserID)

	// Event based
	userLeftEvent := &event.UserLeft{
		ID: userEvent.Payload.UserID,
	}
	uc.room.AddEvent(userLeftEvent)
}

type UserData struct {
	UserID int `json:"user_id"`
	Action struct {
		Name    string `json:"name"`
		Payload struct {
			Axis string `json:"axis"`
			Val  int    `json:"val"`
		} `json:"payload"`
	}
}

func (uc *UserController) FromUser(client bm.Client, data []byte) {
	var userData UserData
	err := json.Unmarshal(data, &userData)
	if err != nil {
		fmt.Println(err)
	}

	if userData.Action.Name == "change_dir" {
		changeDirEvent := &event.ChangeDir{
			ID:   userData.UserID,
			Axis: userData.Action.Payload.Axis,
			Val:  userData.Action.Payload.Val,
		}
		uc.room.AddEvent(changeDirEvent)
	}
}
