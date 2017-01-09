package gate

import (
	"fmt"
	"sync"

	"github.com/jvikstedt/bluemoon/bluemoon"
)

type Hub struct {
	workers     map[int]bluemoon.Client
	users       map[int]bluemoon.Client
	workersLock sync.RWMutex
	usersLock   sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		workers: make(map[int]bluemoon.Client),
		users:   make(map[int]bluemoon.Client),
	}
}

func (h *Hub) AddUser(user bluemoon.Client) error {
	h.usersLock.Lock()
	defer h.usersLock.Unlock()

	if _, ok := h.users[user.ID()]; ok {
		return fmt.Errorf("users map already has a record with id of %d", user.ID())
	}

	h.users[user.ID()] = user

	return nil
}

func (h *Hub) RemoveUser(user bluemoon.Client) error {
	h.usersLock.Lock()
	defer h.usersLock.Unlock()

	if _, ok := h.users[user.ID()]; !ok {
		return fmt.Errorf("user not found with id of %d", user.ID())
	}

	delete(h.users, user.ID())

	return nil
}

func (h *Hub) AddWorker(worker bluemoon.Client) error {
	h.workersLock.Lock()
	defer h.workersLock.Unlock()

	if _, ok := h.workers[worker.ID()]; ok {
		return fmt.Errorf("workers map already has a record with id of %d", worker.ID())
	}

	h.workers[worker.ID()] = worker

	return nil
}

func (h *Hub) RemoveWorker(worker bluemoon.Client) error {
	h.workersLock.Lock()
	defer h.workersLock.Unlock()

	if _, ok := h.workers[worker.ID()]; !ok {
		return fmt.Errorf("worker not found with id of %d", worker.ID())
	}

	delete(h.workers, worker.ID())

	return nil
}

func (h *Hub) UserByID(id int) (bluemoon.Client, error) {
	h.usersLock.RLock()
	defer h.usersLock.RUnlock()

	if v, ok := h.users[id]; ok {
		return v, nil
	}

	return nil, fmt.Errorf("User with id of %d not found", id)
}

func (h *Hub) WorkerByID(id int) (bluemoon.Client, error) {
	h.workersLock.RLock()
	defer h.workersLock.RUnlock()

	if v, ok := h.workers[id]; ok {
		return v, nil
	}

	return nil, fmt.Errorf("Worker with id of %d not found", id)
}
