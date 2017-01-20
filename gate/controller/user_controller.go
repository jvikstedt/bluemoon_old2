package controller

import (
	"encoding/json"
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
	var userIn gate.UserIn
	err := json.Unmarshal(data, &userIn)
	if err != nil {
		fmt.Println(err)
		return
	}

	userInfo, err := uc.uis.ByID(client.ID())
	if err != nil {
		fmt.Println(err)
		return
	}
	worker := userInfo.Worker()

	workerOut := gate.WorkerOut{
		Name: "from_user",
		Payload: gate.Payload{
			Name:    userIn.Name,
			UserID:  client.ID(),
			Payload: userIn.Payload,
		},
	}

	bytes, err := json.Marshal(workerOut)
	if err != nil {
		fmt.Println(err)
		return
	}
	bytes = append(bytes, '\n')

	worker.Write(bytes)
}
