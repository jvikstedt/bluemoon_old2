package bm

import (
	"sync"
)

type UserInfoStore interface {
	ByID(id int) (*UserInfo, error)
	Add(id int, ui *UserInfo)
}

type ClientStore interface {
	Add(user Client) error
	Remove(user Client) error
	ByID(id int) (Client, error)
	One() (Client, error)
}

type UserInfo struct {
	worker Client
	wlock  sync.RWMutex
}

func (uf *UserInfo) Worker() Client {
	uf.wlock.RLock()
	defer uf.wlock.RUnlock()
	return uf.worker
}

func (uf *UserInfo) SetWorker(worker Client) {
	uf.wlock.Lock()
	defer uf.wlock.Unlock()
	uf.worker = worker
}
