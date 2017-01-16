package main

type ChangeDir struct {
	ID   int
	Axis string
	Val  int
}

func (cd *ChangeDir) Execute(room *Room) error {
	entity := room.EntityById(cd.ID)
	playerEntity, ok := entity.(*PlayerEntity)
	if ok {
		if cd.Axis == "x" {
			playerEntity.SetXDir(cd.Val)
		} else if cd.Axis == "y" {
			playerEntity.SetYDir(cd.Val)
		}
	}
	return nil
}
