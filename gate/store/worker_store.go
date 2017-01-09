package store

import (
	"fmt"
	"sync"

	"github.com/jvikstedt/bluemoon/gate/client"
)

type WorkerStore struct {
	workers map[int]*client.Worker
	uLock   sync.RWMutex
}

func NewWorkerStore() *WorkerStore {
	return &WorkerStore{
		workers: make(map[int]*client.Worker),
	}
}

func (us *WorkerStore) Add(worker *client.Worker) error {
	us.uLock.Lock()
	defer us.uLock.Unlock()

	if _, ok := us.workers[worker.ID()]; ok {
		return fmt.Errorf("workers map already has a record with id of %d", worker.ID())
	}

	us.workers[worker.ID()] = worker

	return nil
}

func (us *WorkerStore) Remove(worker *client.Worker) error {
	us.uLock.Lock()
	defer us.uLock.Unlock()

	if _, ok := us.workers[worker.ID()]; !ok {
		return fmt.Errorf("worker not found with id of %d", worker.ID())
	}

	delete(us.workers, worker.ID())

	return nil
}

func (us *WorkerStore) ByID(id int) (*client.Worker, error) {
	us.uLock.RLock()
	defer us.uLock.RUnlock()

	if v, ok := us.workers[id]; ok {
		return v, nil
	}

	return nil, fmt.Errorf("Worker with id of %d not found", id)
}
