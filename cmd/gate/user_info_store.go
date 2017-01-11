package main

import (
	"sync"

	"github.com/jvikstedt/bluemoon/bluemoon"
)

type UserInfoStore struct {
	userInfos map[int]*UserInfo
	uilock    sync.RWMutex
}

func NewUserInfoStore() *UserInfoStore {
	return &UserInfoStore{
		userInfos: make(map[int]*UserInfo),
	}
}

func (uis *UserInfoStore) ByID(id int) *UserInfo {
	uis.uilock.RLock()
	defer uis.uilock.RUnlock()
	return uis.userInfos[id]
}

func (uis *UserInfoStore) Add(id int, ui *UserInfo) {
	uis.uilock.Lock()
	defer uis.uilock.Unlock()
	uis.userInfos[id] = ui
}

type UserInfo struct {
	worker bluemoon.Client
	wlock  sync.RWMutex
}

func (uf *UserInfo) Worker() bluemoon.Client {
	uf.wlock.RLock()
	defer uf.wlock.RUnlock()
	return uf.worker
}

func (uf *UserInfo) SetWorker(worker bluemoon.Client) {
	uf.wlock.Lock()
	defer uf.wlock.Unlock()
	uf.worker = worker
}
