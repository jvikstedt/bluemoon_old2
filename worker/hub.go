package worker

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/jvikstedt/bluemoon/bm"
)

type Hub struct {
	gate  bm.Client
	users map[int]*User
	pLock sync.RWMutex
}

func NewHub(gate bm.Client) *Hub {
	return &Hub{
		gate:  gate,
		users: make(map[int]*User),
	}
}

func (h *Hub) SetGate(gate bm.Client) {
	h.gate = gate
}

type Message struct {
	Name    string      `json:"name"`
	UserIds []int       `json:"user_ids"`
	Payload interface{} `json:"payload"`
}

func (h *Hub) BroadcastTo(userIds []int, payload interface{}) {
	msg := Message{
		Name:    "to_users",
		UserIds: userIds,
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

func (h *Hub) Broadcast(payload interface{}) {
	h.pLock.RLock()
	defer h.pLock.RUnlock()

	ids := make([]int, len(h.users))
	i := 0
	for k := range h.users {
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

func (h *Hub) AddUser(p *User) {
	h.pLock.Lock()
	defer h.pLock.Unlock()
	h.users[p.ID()] = p
}

func (h *Hub) RemoveUser(p *User) {
	h.pLock.Lock()
	defer h.pLock.Unlock()
	delete(h.users, p.ID())
}

func (h *Hub) RemoveUserByID(id int) {
	h.pLock.Lock()
	defer h.pLock.Unlock()
	delete(h.users, id)
}

func (h *Hub) UserByID(id int) *User {
	h.pLock.RLock()
	defer h.pLock.RUnlock()
	return h.users[id]
}
