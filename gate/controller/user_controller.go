package controller

import (
	"fmt"

	"github.com/jvikstedt/bluemoon/bm"
	"github.com/jvikstedt/bluemoon/gate"
)

type UserController struct {
	uis       gate.UserInfoStore
	userStore gate.ClientStore
}

func NewUserController(uis gate.UserInfoStore, userStore gate.ClientStore) *UserController {
	return &UserController{
		uis:       uis,
		userStore: userStore,
	}
}

func (uc *UserController) ToWorker(client bm.Client, data []byte) {
	userInfo, err := uc.uis.ByID(client.ID())
	if err != nil {
		fmt.Println(err)
		return
	}
	worker := userInfo.Worker()
	data = append(data, '\n')
	worker.Write(data)
}
