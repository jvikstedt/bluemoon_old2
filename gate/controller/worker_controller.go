package controller

import "github.com/jvikstedt/bluemoon/gate"

type WorkerController struct {
	uis       gate.UserInfoStore
	userStore gate.ClientStore
}

func NewWorkerController(uis gate.UserInfoStore, userStore gate.ClientStore) *WorkerController {
	return &WorkerController{
		uis:       uis,
		userStore: userStore,
	}
}
