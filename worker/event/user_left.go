package event

import (
	"fmt"

	"github.com/jvikstedt/bluemoon/worker/room"
)

type UserLeft struct {
	ID int
}

func (uj *UserLeft) Execute(r *room.Room) error {
	r.RemoveUserByID(uj.ID)
	r.RemoveEntityByID(uj.ID)
	r.Broadcast([]byte(fmt.Sprintf(`{"name": "remove_player", "id": %d}`, uj.ID)))
	return nil
}
