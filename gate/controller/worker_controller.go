package controller

import (
	"encoding/json"
	"fmt"

	"github.com/jvikstedt/bluemoon/bm"
	"github.com/jvikstedt/bluemoon/gate"
)

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

type WorkerData struct {
	UserIds []int            `json:"user_ids"`
	Payload *json.RawMessage `json:"payload"`
}

func (wc *WorkerController) ToUsers(client bm.Client, data []byte) {
	var workerData WorkerData
	err := json.Unmarshal(data, &workerData)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, userID := range workerData.UserIds {
		user, err := wc.userStore.ByID(userID)
		if err != nil {
			fmt.Println(err)
			return
		}
		user.Write(*workerData.Payload)
	}
}
