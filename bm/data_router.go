package bm

import (
	"fmt"
	"sync"
)

type DataRouter struct {
	handlers     map[string]HandleClientDataFunc
	defaultRoute HandleClientDataFunc
	hlock        sync.RWMutex
}

func NewDataRouter() *DataRouter {
	return &DataRouter{
		handlers: make(map[string]HandleClientDataFunc),
	}
}

func (dr *DataRouter) Register(name string, handler HandleClientDataFunc) {
	dr.hlock.Lock()
	defer dr.hlock.Unlock()

	dr.handlers[name] = handler
}

func (dr *DataRouter) UnRegister(name string) error {
	dr.hlock.Lock()
	defer dr.hlock.Unlock()

	if _, ok := dr.handlers[name]; ok {
		delete(dr.handlers, name)
		return nil
	}
	return fmt.Errorf("Could not find handler with a name: %s", name)
}

func (dr *DataRouter) Handler(name string) (HandleClientDataFunc, error) {
	dr.hlock.RLock()
	defer dr.hlock.RUnlock()

	if h, ok := dr.handlers[name]; ok {
		return h, nil
	}
	if dr.defaultRoute != nil {
		return dr.defaultRoute, nil
	}
	return nil, fmt.Errorf("Could not find handler with a name: %s", name)
}

func (dr *DataRouter) SetDefaultHandler(handler HandleClientDataFunc) {
	dr.hlock.Lock()
	defer dr.hlock.Unlock()

	dr.defaultRoute = handler
}
