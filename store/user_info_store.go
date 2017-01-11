package store

import (
	"sync"

	"github.com/jvikstedt/bluemoon/bluemoon"
)

type UserInfoStore struct {
	userInfos map[int]*bluemoon.UserInfo
	uilock    sync.RWMutex
}

func NewUserInfoStore() *UserInfoStore {
	return &UserInfoStore{
		userInfos: make(map[int]*bluemoon.UserInfo),
	}
}

func (uis *UserInfoStore) ByID(id int) *bluemoon.UserInfo {
	uis.uilock.RLock()
	defer uis.uilock.RUnlock()
	return uis.userInfos[id]
}

func (uis *UserInfoStore) Add(id int, ui *bluemoon.UserInfo) {
	uis.uilock.Lock()
	defer uis.uilock.Unlock()
	uis.userInfos[id] = ui
}
