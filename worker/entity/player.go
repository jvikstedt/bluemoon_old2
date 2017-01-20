package entity

type PlayerEntity struct {
	id    int
	x     int
	y     int
	xDir  int
	yDir  int
	speed float64
}

func NewPlayerEntity(id, x, y, xdir, ydir int, speed float64) *PlayerEntity {
	return &PlayerEntity{
		id:    id,
		x:     x,
		y:     y,
		xDir:  xdir,
		yDir:  ydir,
		speed: speed,
	}
}

func (pe *PlayerEntity) Update(delta float64) (changed bool, err error) {
	oldX := pe.x
	oldY := pe.y

	pe.x = pe.x + int(float64(pe.xDir)*pe.speed*delta)
	pe.y = pe.y + int(float64(pe.yDir)*pe.speed*delta)

	if oldX != pe.x || oldY != pe.y {
		changed = true
	}
	return
}

func (pe *PlayerEntity) ID() int {
	return pe.id
}

func (pe *PlayerEntity) X() int {
	return pe.x
}

func (pe *PlayerEntity) Y() int {
	return pe.y
}

func (pe *PlayerEntity) SetXDir(dir int) {
	pe.xDir = dir
}

func (pe *PlayerEntity) SetYDir(dir int) {
	pe.yDir = dir
}
