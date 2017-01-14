package main

import (
	"sync"

	"github.com/jvikstedt/bluemoon/bm"
)

type Player struct {
	client bm.Client
	id     int
	x      int
	y      int
	pLock  sync.RWMutex
}

func NewPlayer(client bm.Client, id, x, y int) *Player {
	return &Player{
		client: client,
		id:     id,
		x:      x,
		y:      y,
	}
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

func (p *Player) Client() bm.Client {
	return p.client
}
