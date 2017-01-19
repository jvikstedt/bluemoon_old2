package event

import (
	"fmt"

	"github.com/jvikstedt/bluemoon/worker"
	"github.com/jvikstedt/bluemoon/worker/entity"
	"github.com/jvikstedt/bluemoon/worker/room"
)

type UserJoined struct {
	ID int
}

func (uj *UserJoined) Execute(r worker.Room) error {
	game, ok := r.(*room.Game)
	if !ok {
		return fmt.Errorf("Wrong type of room, expected Game")
	}
	game.AddUser(worker.NewUser(uj.ID))
	game.AddEntity(entity.NewPlayerEntity(uj.ID, 50, 50, 0, 0, 4))
	for _, v := range game.Entities() {
		game.BroadcastTo([]int{uj.ID}, struct {
			Name string `json:"name"`
			ID   int    `json:"id"`
			X    int    `json:"x"`
			Y    int    `json:"y"`
		}{"new_player", v.ID(), v.X(), v.Y()})
	}
	return nil
}
