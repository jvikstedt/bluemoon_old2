package gate

import (
	"sync"

	"github.com/jvikstedt/bluemoon/bm"
)

type UserInfoStore interface {
	ByID(id int) (*UserInfo, error)
	Add(id int, ui *UserInfo)
}

type ClientStore interface {
	Add(user bm.Client) error
	Remove(user bm.Client) error
	ByID(id int) (bm.Client, error)
	One() (bm.Client, error)
}

type UserInfo struct {
	worker bm.Client
	wlock  sync.RWMutex
}

func (uf *UserInfo) Worker() bm.Client {
	uf.wlock.RLock()
	defer uf.wlock.RUnlock()
	return uf.worker
}

func (uf *UserInfo) SetWorker(worker bm.Client) {
	uf.wlock.Lock()
	defer uf.wlock.Unlock()
	uf.worker = worker
}
