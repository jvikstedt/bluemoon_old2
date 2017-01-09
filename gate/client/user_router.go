package client

import (
	"fmt"
	"sync"
)

type UserRouter struct {
	handlers map[string]HandleUserDataFunc
	hlock    sync.RWMutex
}

func NewUserRouter() *UserRouter {
	return &UserRouter{
		handlers: make(map[string]HandleUserDataFunc),
	}
}

func (dr *UserRouter) Register(name string, handler HandleUserDataFunc) {
	dr.hlock.Lock()
	defer dr.hlock.Unlock()

	dr.handlers[name] = handler
}

func (dr *UserRouter) UnRegister(name string) error {
	dr.hlock.Lock()
	defer dr.hlock.Unlock()

	if _, ok := dr.handlers[name]; ok {
		delete(dr.handlers, name)
		return nil
	}
	return fmt.Errorf("Could not find handler with a name: %s", name)
}

func (dr *UserRouter) Handler(name string) (HandleUserDataFunc, error) {
	dr.hlock.RLock()
	defer dr.hlock.RUnlock()

	if h, ok := dr.handlers[name]; ok {
		return h, nil
	}
	return nil, fmt.Errorf("Could not find handler with a name: %s", name)
}
