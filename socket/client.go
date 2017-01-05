package socket

import (
	"fmt"
	"net"
	"time"

	"github.com/jvikstedt/bluemoon/gate"
)

type Client struct {
}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) Connect(addr string) error {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", addr)
	if err != nil {
		return err
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return err
	}

	cw := NewConnectionWrapper(conn)
	defer cw.Close()

	worker := gate.NewWorker(1, cw, func(worker *gate.Worker, data []byte) {
		fmt.Printf("message %s\n", string(data))
	})

	go worker.EnableReader()
	go worker.EnableWriter()

	for {
		worker.Write([]byte("ping"))
		time.Sleep(time.Second * 1)
	}

	return nil
}
