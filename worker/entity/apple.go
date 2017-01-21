package entity

type AppleEntity struct {
	id int
	x  int
	y  int
}

func NewAppleEntity(id, x, y int) *AppleEntity {
	return &AppleEntity{
		id: id,
		x:  x,
		y:  y,
	}
}

func (ae *AppleEntity) Update(delta float64) (changed bool, err error) {
	return
}

func (pe *AppleEntity) ID() int {
	return pe.id
}

func (pe *AppleEntity) X() int {
	return pe.x
}

func (pe *AppleEntity) Y() int {
	return pe.y
}

func (pe *AppleEntity) Type() string {
	return "apple"
}
