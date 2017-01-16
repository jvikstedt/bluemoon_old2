package main

import "sync"

type Player struct {
	id    int
	pLock sync.RWMutex
}

func NewPlayer(id int) *Player {
	return &Player{
		id: id,
	}
}

func (p *Player) ID() int {
	p.pLock.RLock()
	defer p.pLock.RUnlock()
	return p.id
}
