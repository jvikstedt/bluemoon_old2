package main

import (
	"fmt"
	"sync"
)

type Hub struct {
	players map[int]*Player
	pLock   sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		players: make(map[int]*Player),
	}
}

func (h *Hub) AddPlayer(p *Player) {
	h.pLock.Lock()
	defer h.pLock.Unlock()
	h.players[p.ID()] = p
}

func (h *Hub) RemovePlayer(p *Player) {
	h.pLock.Lock()
	defer h.pLock.Unlock()
	delete(h.players, p.ID())
}

func (h *Hub) RemovePlayerByID(id int) {
	h.pLock.Lock()
	defer h.pLock.Unlock()
	delete(h.players, id)
}

func (h *Hub) PlayerByID(id int) *Player {
	h.pLock.RLock()
	defer h.pLock.RUnlock()
	fmt.Println(h.players)
	return h.players[id]
}
