package event

import (
	"github.com/jvikstedt/bluemoon/worker/entity"
	"github.com/jvikstedt/bluemoon/worker/room"
)

type ChangeDir struct {
	ID   int
	Axis string
	Val  int
}

func (cd *ChangeDir) Execute(room *room.Room) error {
	e := room.EntityById(cd.ID)
	playerEntity, ok := e.(*entity.PlayerEntity)
	if ok {
		if cd.Axis == "x" {
			playerEntity.SetXDir(cd.Val)
		} else if cd.Axis == "y" {
			playerEntity.SetYDir(cd.Val)
		}
	}
	return nil
}
