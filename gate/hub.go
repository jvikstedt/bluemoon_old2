package gate

import (
	"fmt"
	"sync"
)

type Hub struct {
	workers     map[int]*Worker
	users       map[int]*User
	workersLock sync.RWMutex
	usersLock   sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		workers: make(map[int]*Worker),
		users:   make(map[int]*User),
	}
}

func (h *Hub) UserByID(id int) (*User, error) {
	h.usersLock.RLock()
	defer h.usersLock.RUnlock()

	if v, ok := h.users[id]; ok {
		return v, nil
	}

	return nil, fmt.Errorf("User with id of %d not found", id)
}

func (h *Hub) WorkerByID(id int) (*Worker, error) {
	h.workersLock.RLock()
	defer h.workersLock.RUnlock()

	if v, ok := h.workers[id]; ok {
		return v, nil
	}

	return nil, fmt.Errorf("Worker with id of %d not found", id)
}
