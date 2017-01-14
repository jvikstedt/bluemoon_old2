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
}

func (uc *UserController) UserLeft(client bm.Client, data []byte) {
	var userEvent UserEvent
	err := json.Unmarshal(data, &userEvent)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(userEvent)
	uc.hub.RemovePlayerByID(client.ID())
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
		x := player.X()
		x += 20 * moveEvent.Payload.Val
		player.SetX(x)
	} else {
		y := player.Y()
		y += 20 * moveEvent.Payload.Val
		player.SetY(y)
	}

	client.Write([]byte(fmt.Sprintf(`{"name": "move", "user_id": %d, "payload": {"id": %d, "x": %d, "y": %d}}`, moveEvent.UserID, moveEvent.UserID, player.X(), player.Y()) + "\n"))
}
