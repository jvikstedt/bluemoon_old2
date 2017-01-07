package gate

import (
	"fmt"
	"sync"
)

type Hub struct {
	workers     map[int]*Worker
	users       map[int]*User
	workersLock sync.Mutex
	usersLock   sync.Mutex
}

func NewHub() *Hub {
	return &Hub{
		workers: make(map[int]*Worker),
		users:   make(map[int]*User),
	}
}

func (h *Hub) UserByID(id int) (*User, error) {
	h.usersLock.Lock()
	defer h.usersLock.Unlock()

	if v, ok := h.users[id]; ok {
		return v, nil
	}

	return nil, fmt.Errorf("User with id of %d not found", id)
}

func (h *Hub) WorkerByID(id int) (*Worker, error) {
	h.workersLock.Lock()
	defer h.workersLock.Unlock()

	if v, ok := h.workers[id]; ok {
		return v, nil
	}

	return nil, fmt.Errorf("Worker with id of %d not found", id)
}
