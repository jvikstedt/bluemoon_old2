package ws

import "github.com/gorilla/websocket"

type Client struct {
}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) Connect(addr string) (*websocket.Conn, error) {
	conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:4000/ws", nil)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
