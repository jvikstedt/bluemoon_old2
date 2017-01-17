package event

import (
	"fmt"

	"github.com/jvikstedt/bluemoon/worker/room"
)

type UserLeft struct {
	ID int
}

func (uj *UserLeft) Execute(room *room.Room) error {
	entity := room.EntityById(uj.ID)
	room.RemoveEntity(entity)
	room.Hub().Broadcast([]byte(fmt.Sprintf(`{"name": "remove_player", "id": %d}`, uj.ID)))
	return nil
}
