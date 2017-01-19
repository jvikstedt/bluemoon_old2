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

func (h *Hub) ManageWorkerConn(rw bm.ReadWriter) error {
	w := bm.NewBaseClient(idgen.Next(), rw, func(client bm.Client, data []byte) {
		fmt.Printf("New message from worker: %d\n", client.ID())
		fmt.Print(string(data))
		var workerIn WorkerIn
		err := json.Unmarshal(data, &workerIn)
		if err != nil {
			fmt.Println(err)
			return
		}
		handle, err := h.workerRouter.Handler(workerIn.Name)
		if err != nil {
			fmt.Println(err)
			return
		}
		handle(client, *workerIn.Payload)
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

func (h *Hub) ManageUserConn(rw bm.ReadWriter) error {
	u := bm.NewBaseClient(idgen.Next(), rw, func(client bm.Client, data []byte) {
		fmt.Printf("New message from user: %d\n", client.ID())
		fmt.Print(string(data))
		var userIn UserIn
		err := json.Unmarshal(data, &userIn)
		if err != nil {
			fmt.Println(err)
			return
		}
		handle, err := h.userRouter.Handler(userIn.Name)
		handle(client, data)
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
