package game

import (
	"io"
	"sync"

	"github.com/jvikstedt/bluemoon/log"
)

type Room interface {
	Type() string
	ID() int
}

type packer interface {
	Unpack(data []byte) (name string, payload []byte, err error)
	Pack(name string, payload []byte) (data []byte, err error)
}

type Connection io.ReadWriter

type Hub struct {
	log   log.Logger
	mu    sync.Mutex
	users map[int]User
	p     packer
}

func NewHub(log log.Logger, p packer) *Hub {
	return &Hub{
		log:   log,
		users: make(map[int]User),
		p:     p,
	}
}

func (h *Hub) HandleConnection(conn Connection) {
	client := NewClient(h.log, 1, conn, h.handleData)
	go client.EnableReader()
	client.EnableWriter()
}

func (h *Hub) handleData(c *Client, bytes []byte) {
	h.log.Printf("Client %d sent: %s", c.ID(), string(bytes))
	name, payload, err := h.p.Unpack(bytes)
	if err != nil {
		h.log.Printf("%v", err)
	}
	h.log.Printf("Got %s and %s\n", name, string(payload))
}
