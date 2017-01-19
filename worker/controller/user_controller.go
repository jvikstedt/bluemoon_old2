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

type ChangeDirAction struct {
	UserID  int `json:"user_id"`
	Payload struct {
		Axis string `json:"axis"`
		Val  int    `json:"val"`
	} `json:"payload"`
}

func (uc *UserController) ChangeDir(client bm.Client, data []byte) {
	var cd ChangeDirAction
	err := json.Unmarshal(data, &cd)
	if err != nil {
		fmt.Println(err)
		return
	}

	changeDirEvent := &event.ChangeDir{
		ID:   cd.UserID,
		Axis: cd.Payload.Axis,
		Val:  cd.Payload.Val,
	}
	uc.room.AddEvent(changeDirEvent)
}
