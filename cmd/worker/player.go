package main

import "sync"

type Player struct {
	id    int
	x     int
	y     int
	xDir  int
	yDir  int
	speed float64
	pLock sync.RWMutex
}

func NewPlayer(id, x, y int) *Player {
	return &Player{
		id:    id,
		x:     x,
		y:     y,
		speed: 10,
	}
}

func (p *Player) Update(delta float64) {
	p.pLock.Lock()
	defer p.pLock.Unlock()

	p.x = p.x + int(float64(p.xDir)*p.speed*delta)
	p.y = p.y + int(float64(p.yDir)*p.speed*delta)
}

func (p *Player) ID() int {
	p.pLock.RLock()
	defer p.pLock.RUnlock()
	return p.id
}

func (p *Player) X() int {
	p.pLock.RLock()
	defer p.pLock.RUnlock()
	return p.x
}

func (p *Player) Y() int {
	p.pLock.RLock()
	defer p.pLock.RUnlock()
	return p.y
}

func (p *Player) SetID(id int) {
	p.pLock.Lock()
	defer p.pLock.Unlock()
	p.id = id
}

func (p *Player) SetY(y int) {
	p.pLock.Lock()
	defer p.pLock.Unlock()
	p.y = y
}

func (p *Player) SetX(x int) {
	p.pLock.Lock()
	defer p.pLock.Unlock()
	p.x = x
}

func (p *Player) SetXDir(dir int) {
	p.pLock.Lock()
	defer p.pLock.Unlock()
	p.xDir = dir
}

func (p *Player) SetYDir(dir int) {
	p.pLock.Lock()
	defer p.pLock.Unlock()
	p.yDir = dir
}
