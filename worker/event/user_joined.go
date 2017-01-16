package event

import (
	"fmt"

	"github.com/jvikstedt/bluemoon/worker/entity"
	"github.com/jvikstedt/bluemoon/worker/room"
)

type UserJoined struct {
	ID int
}

func (uj *UserJoined) Execute(room *room.Room) error {
	room.AddEntity(entity.NewPlayerEntity(uj.ID, 50, 50, 0, 0, 4))
	for _, v := range room.Entities() {
		room.Hub().BroadcastTo([]int{uj.ID}, []byte(fmt.Sprintf(`{"name": "new_player", "id": %d, "x": %d, "y": %d}`, v.ID(), v.X(), v.Y())))
	}
	return nil
}
