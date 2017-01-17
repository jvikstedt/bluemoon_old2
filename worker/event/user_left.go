package event

import (
	"fmt"

	"github.com/jvikstedt/bluemoon/worker"
	"github.com/jvikstedt/bluemoon/worker/room"
)

type UserLeft struct {
	ID int
}

func (uj *UserLeft) Execute(r worker.Room) error {
	game, ok := r.(*room.Game)
	if !ok {
		return fmt.Errorf("Wrong type of room, expected Game")
	}
	game.RemoveUserByID(uj.ID)
	game.RemoveEntityByID(uj.ID)
	game.Broadcast([]byte(fmt.Sprintf(`{"name": "remove_player", "id": %d}`, uj.ID)))
	return nil
}
