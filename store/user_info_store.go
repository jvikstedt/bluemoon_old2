package store

import (
	"fmt"
	"sync"

	"github.com/jvikstedt/bluemoon/bm"
)

type UserInfoStore struct {
	userInfos map[int]*bm.UserInfo
	uilock    sync.RWMutex
}

func NewUserInfoStore() *UserInfoStore {
	return &UserInfoStore{
		userInfos: make(map[int]*bm.UserInfo),
	}
}

func (uis *UserInfoStore) ByID(id int) (*bm.UserInfo, error) {
	uis.uilock.RLock()
	defer uis.uilock.RUnlock()
	if ui, ok := uis.userInfos[id]; ok {
		return ui, nil
	}
	return nil, fmt.Errorf("No UserInfo found with an id of: %d", id)
}

func (uis *UserInfoStore) Add(id int, ui *bm.UserInfo) {
	uis.uilock.Lock()
	defer uis.uilock.Unlock()
	uis.userInfos[id] = ui
}
