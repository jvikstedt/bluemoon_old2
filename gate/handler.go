package gate

import "fmt"

type Handler struct {
	hub *Hub
}

func NewHandler(hub *Hub) *Handler {
	return &Handler{
		hub: hub,
	}
}

func (h *Handler) HandleUserData(user *User, data []byte) {
	fmt.Printf("Got message from user: %s\n", string(data))
	user.Write([]byte("Pong"))
}

func (h *Handler) HandleWorkerData(worker *Worker, data []byte) {
	fmt.Printf("Got message from user: %s\n", string(data))
	worker.Write([]byte("Pong"))
}
