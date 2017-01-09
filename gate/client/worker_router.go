package client

import (
	"fmt"
	"sync"
)

type WorkerRouter struct {
	handlers map[string]HandleWorkerDataFunc
	hlock    sync.RWMutex
}

func NewWorkerRouter() *WorkerRouter {
	return &WorkerRouter{
		handlers: make(map[string]HandleWorkerDataFunc),
	}
}

func (dr *WorkerRouter) Register(name string, handler HandleWorkerDataFunc) {
	dr.hlock.Lock()
	defer dr.hlock.Unlock()

	dr.handlers[name] = handler
}

func (dr *WorkerRouter) UnRegister(name string) error {
	dr.hlock.Lock()
	defer dr.hlock.Unlock()

	if _, ok := dr.handlers[name]; ok {
		delete(dr.handlers, name)
		return nil
	}
	return fmt.Errorf("Could not find handler with a name: %s", name)
}

func (dr *WorkerRouter) Handler(name string) (HandleWorkerDataFunc, error) {
	dr.hlock.RLock()
	defer dr.hlock.RUnlock()

	if h, ok := dr.handlers[name]; ok {
		return h, nil
	}
	return nil, fmt.Errorf("Could not find handler with a name: %s", name)
}
