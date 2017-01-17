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
	uc.hub.Broadcast([]byte(fmt.Sprintf(`{"name": "new_player", "id": %d, "x": 50, "y": 50}`, user.ID())))
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

type MoveEvent struct {
	Name    string `json:"name"`
	UserID  int    `json:"user_id"`
	Payload struct {
		Axis string `json:"axis"`
		Val  int    `json:"val"`
	} `json:"payload"`
}

func (uc *UserController) Direction(client bm.Client, data []byte) {
	var moveEvent MoveEvent
	err := json.Unmarshal(data, &moveEvent)
	if err != nil {
		fmt.Println(err)
	}

	// Event based
	changeDirEvent := &event.ChangeDir{
		ID:   moveEvent.UserID,
		Axis: moveEvent.Payload.Axis,
		Val:  moveEvent.Payload.Val,
	}
	uc.room.AddEvent(changeDirEvent)
}
