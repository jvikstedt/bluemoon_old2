package main

import (
	"encoding/json"
	"fmt"

	"github.com/jvikstedt/bluemoon/bm"
)

type Hub struct {
	dataRouter    *bm.DataRouter
	workerStore   bm.ClientStore
	userStore     bm.ClientStore
	userInfoStore bm.UserInfoStore
}

func NewHub(dr *bm.DataRouter, ws bm.ClientStore, us bm.ClientStore, uis bm.UserInfoStore) *Hub {
	return &Hub{
		dataRouter:    dr,
		workerStore:   ws,
		userStore:     us,
		userInfoStore: uis,
	}
}

type DN struct {
	Name string `json:"name"`
}

var idgen = bm.NewIDGen()

func (h *Hub) ManageWorkerConn(rw bm.ReadWriter) error {
	w := bm.NewBaseClient(idgen.Next(), rw, func(client bm.Client, data []byte) {
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

func (h *Hub) buildUserInfo(client bm.Client, worker bm.Client) {
	ui := &bm.UserInfo{}
	ui.SetWorker(worker)
	h.userInfoStore.Add(client.ID(), ui)
}

func (h *Hub) ManageUserConn(rw bm.ReadWriter) error {
	u := bm.NewBaseClient(idgen.Next(), rw, func(client bm.Client, data []byte) {
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

	worker, err := h.workerStore.One()
	if err != nil {
		return err
	}
	h.buildUserInfo(u, worker)

	h.userStore.Add(u)
	defer h.userStore.Remove(u)

	go u.EnableReader()
	u.EnableWriter()

	return nil
}
