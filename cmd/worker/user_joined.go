package main

import "fmt"

type PlayerEntity struct {
	id    int
	x     int
	y     int
	xDir  int
	yDir  int
	speed float64
}

func (pe *PlayerEntity) Update(delta float64) (bool, error) {
	pe.x = pe.x + int(float64(pe.xDir)*pe.speed*delta)
	pe.y = pe.y + int(float64(pe.yDir)*pe.speed*delta)
	return true, nil
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

type UserJoined struct {
	ID int
}

func (uj *UserJoined) Execute(room *Room) error {
	room.AddEntity(&PlayerEntity{id: uj.ID, x: 50, y: 50, speed: 4})
	for _, v := range room.entities {
		room.hub.BroadcastTo([]int{uj.ID}, []byte(fmt.Sprintf(`{"name": "new_player", "id": %d, "x": %d, "y": %d}`, v.ID(), v.X(), v.Y())))
	}
	return nil
}
