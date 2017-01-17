package event

import (
	"fmt"
	"github.com/jvikstedt/bluemoon/worker"
	"github.com/jvikstedt/bluemoon/worker/entity"
)

type ChangeDir struct {
	ID   int
	Axis string
	Val  int
}

func (cd *ChangeDir) Execute(room *worker.Room) error {
	e := room.EntityById(cd.ID)
	playerEntity, ok := e.(*entity.PlayerEntity)
	if !ok {
		return fmt.Errorf("Wrong type of entity, expected PlayerEntity with id of: %d", cd.ID)
	}

	if cd.Axis == "x" {
		playerEntity.SetXDir(cd.Val)
	} else if cd.Axis == "y" {
		playerEntity.SetYDir(cd.Val)
	}
	return nil
}
