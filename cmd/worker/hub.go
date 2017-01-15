package main

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/jvikstedt/bluemoon/bm"
)

type Hub struct {
	gate    bm.Client
	players map[int]*Player
	pLock   sync.RWMutex
}

func NewHub(gate bm.Client) *Hub {
	return &Hub{
		gate:    gate,
		players: make(map[int]*Player),
	}
}

func (h *Hub) SetGate(gate bm.Client) {
	h.gate = gate
}

func (h *Hub) Run() {
	last := time.Now()
	for range time.Tick(1 * time.Second) {
		delta := time.Since(last).Seconds() + 1
		last = time.Now()

		h.updatePlayers(delta)
	}
}

func (h *Hub) updatePlayers(delta float64) {
	h.pLock.RLock()
	defer h.pLock.RUnlock()

	for _, p := range h.players {
		p.Update(delta)
		h.Broadcast([]byte(fmt.Sprintf(`{"name": "move", "id": %d, "x": %d, "y": %d}`, p.ID(), p.X(), p.Y())))
	}
}

type Message struct {
	Name    string `json:"name"`
	UserIds []int  `json:"user_ids"`
	Payload []byte `json:"payload"`
}

func (h *Hub) Broadcast(payload []byte) {
	h.pLock.RLock()
	defer h.pLock.RUnlock()

	ids := make([]int, len(h.players))
	i := 0
	for k := range h.players {
		ids[i] = k
		i++
	}

	msg := Message{
		Name:    "to_users",
		UserIds: ids,
		Payload: payload,
	}
	bytes, err := json.Marshal(msg)
	if err != nil {
		fmt.Println(err)
		return
	}

	bytes = append(bytes, '\n')
	h.gate.Write(bytes)
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
