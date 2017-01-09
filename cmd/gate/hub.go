package main

import (
	"encoding/json"
	"fmt"

	"github.com/jvikstedt/bluemoon/bluemoon"
)

type Hub struct {
	dataRouter  *bluemoon.DataRouter
	workerStore *bluemoon.ClientStore
	userStore   *bluemoon.ClientStore
}

func NewHub(dr *bluemoon.DataRouter, ws *bluemoon.ClientStore, us *bluemoon.ClientStore) *Hub {
	return &Hub{
		dataRouter:  dr,
		workerStore: ws,
		userStore:   us,
	}
}

type DN struct {
	Name string `json:"name"`
}

func (h *Hub) ManageWorkerConn(rw bluemoon.ReadWriter) error {
	w := bluemoon.NewBaseClient(1, rw, func(client bluemoon.Client, data []byte) {
		fmt.Printf("New message from worker: %d\n", client.ID())
		fmt.Print(string(data))
		var dn DN
		err := json.Unmarshal(data, &dn)
		if err != nil {
			fmt.Println(err)
			return
		}
		handle, err := h.dataRouter.Handler(dn.Name)
		if err != nil {
			fmt.Println(err)
			return
		}
		handle(client, data)
	})
	h.workerStore.Add(w)
	defer h.workerStore.Remove(w)

	go w.EnableReader()
	w.EnableWriter()

	return nil
}

func (h *Hub) ManageUserConn(rw bluemoon.ReadWriter) error {
	u := bluemoon.NewBaseClient(1, rw, func(client bluemoon.Client, data []byte) {
		fmt.Printf("New message from user: %d\n", client.ID())
		fmt.Print(string(data))
		var dn DN
		err := json.Unmarshal(data, &dn)
		if err != nil {
			fmt.Println(err)
			return
		}
		handle, err := h.dataRouter.Handler(dn.Name)
		if err != nil {
			fmt.Println(err)
			return
		}
		handle(client, data)
	})

	h.userStore.Add(u)
	defer h.userStore.Remove(u)

	go u.EnableReader()
	u.EnableWriter()

	return nil
}
