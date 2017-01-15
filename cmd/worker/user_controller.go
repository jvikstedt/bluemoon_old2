package main

import (
	"encoding/json"
	"fmt"

	"github.com/jvikstedt/bluemoon/bm"
)

type UserController struct {
	hub *Hub
}

func NewUserController(hub *Hub) *UserController {
	return &UserController{
		hub: hub,
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
	player := NewPlayer(userEvent.Payload.UserID, 50, 50)
	uc.hub.AddPlayer(player)

	uc.hub.Broadcast([]byte(fmt.Sprintf(`{"name": "new_player", "id": %d, "x": 50, "y": 50}`, player.ID())))
}

func (uc *UserController) UserLeft(client bm.Client, data []byte) {
	var userEvent UserEvent
	err := json.Unmarshal(data, &userEvent)
	if err != nil {
		fmt.Println(err)
	}
	uc.hub.RemovePlayerByID(userEvent.Payload.UserID)
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
	player := uc.hub.PlayerByID(moveEvent.UserID)
	if player == nil {
		fmt.Printf("Player nil: %d\n", moveEvent.UserID)
		return
	}

	if moveEvent.Payload.Axis == "x" {
		player.SetXDir(moveEvent.Payload.Val)
	} else if moveEvent.Payload.Axis == "y" {
		player.SetYDir(moveEvent.Payload.Val)
	}
}
