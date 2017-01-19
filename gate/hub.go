package gate

import (
	"encoding/json"
	"fmt"

	"github.com/jvikstedt/bluemoon/bm"
)

type Hub struct {
	userRouter    *bm.DataRouter
	workerRouter  *bm.DataRouter
	workerStore   ClientStore
	userStore     ClientStore
	userInfoStore UserInfoStore
}

func NewHub(userRouter *bm.DataRouter, workerRouter *bm.DataRouter, ws ClientStore, us ClientStore, uis UserInfoStore) *Hub {
	return &Hub{
		userRouter:    userRouter,
		workerRouter:  workerRouter,
		workerStore:   ws,
		userStore:     us,
		userInfoStore: uis,
	}
}

var idgen = bm.NewIDGen()

type WorkerData struct {
	Name    string           `json:"name"`
	UserIds []int            `json:"user_ids"`
	Payload *json.RawMessage `json:"payload"`
}

func (h *Hub) ManageWorkerConn(rw bm.ReadWriter) error {
	w := bm.NewBaseClient(idgen.Next(), rw, func(client bm.Client, data []byte) {
		fmt.Printf("New message from worker: %d\n", client.ID())
		fmt.Print(string(data))
		var workerData WorkerData
		err := json.Unmarshal(data, &workerData)
		if err != nil {
			fmt.Println(err)
			return
		}
		if handle, err := h.workerRouter.Handler(workerData.Name); err == nil {
			handle(client, data)
		} else {
			for _, userID := range workerData.UserIds {
				user, err := h.userStore.ByID(userID)
				if err != nil {
					fmt.Println(err)
					return
				}
				user.Write(*workerData.Payload)
			}
		}
	})
	h.workerStore.Add(w)
	defer h.workerStore.Remove(w)

	go w.EnableReader()
	w.EnableWriter()

	return nil
}

func (h *Hub) buildUserInfo(client bm.Client, worker bm.Client) {
	ui := &UserInfo{}
	ui.SetWorker(worker)
	h.userInfoStore.Add(client.ID(), ui)
}

func (h *Hub) userQuit(client bm.Client) {
	userInfo, err := h.userInfoStore.ByID(client.ID())
	if err != nil {
		return
	}
	worker := userInfo.Worker()
	if worker != nil {
		worker.Write([]byte(fmt.Sprintf(`{"name": "user_left", "payload": {"user_id": %d}}`, client.ID()) + "\n"))
	}
}

type UserData struct {
	Name    string           `json:"name"`
	UserID  int              `json:"user_id"`
	Payload *json.RawMessage `json:"payload"`
}

type ToWorkerData struct {
	Name   string `json:"name"`
	UserID int    `json:"user_id"`
	Action struct {
		Name    string           `json:"name"`
		Payload *json.RawMessage `json:"payload"`
	}
}

func (h *Hub) ManageUserConn(rw bm.ReadWriter) error {
	u := bm.NewBaseClient(idgen.Next(), rw, func(client bm.Client, data []byte) {
		fmt.Printf("New message from user: %d\n", client.ID())
		fmt.Print(string(data))
		var userData UserData
		err := json.Unmarshal(data, &userData)
		if err != nil {
			fmt.Println(err)
			return
		}
		if handle, err := h.userRouter.Handler(userData.Name); err == nil {
			handle(client, data)
		} else {
			handle, err = h.userRouter.Handler("ToWorker")
			if err != nil {
				fmt.Println(err)
				return
			}
			var toWorker ToWorkerData
			toWorker.Name = "FromUser"
			toWorker.UserID = client.ID()
			toWorker.Action.Name = userData.Name
			toWorker.Action.Payload = userData.Payload
			bytes, err := json.Marshal(toWorker)
			if err != nil {
				fmt.Println(err)
				return
			}
			handle(client, bytes)
		}
	})

	worker, err := h.workerStore.One()
	if err != nil {
		return err
	}

	worker.Write([]byte(fmt.Sprintf(`{"name": "user_joined", "payload": {"user_id": %d}}`, u.ID()) + "\n"))
	h.buildUserInfo(u, worker)
	defer h.userQuit(u)

	h.userStore.Add(u)
	defer h.userStore.Remove(u)

	go u.EnableReader()
	u.EnableWriter()

	return nil
}
