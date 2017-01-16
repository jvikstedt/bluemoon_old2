package ws

import (
	"github.com/gorilla/websocket"
)

type ConnectionWrapper struct {
	conn *websocket.Conn
}

func NewConnectionWrapper(conn *websocket.Conn) *ConnectionWrapper {
	return &ConnectionWrapper{
		conn: conn,
	}
}

func (cw *ConnectionWrapper) Close() error {
	return cw.conn.Close()
}

func (cw *ConnectionWrapper) Write(data []byte) {
	cw.conn.WriteMessage(websocket.TextMessage, data)
}

func (cw *ConnectionWrapper) Read() ([]byte, error) {
	_, data, err := cw.conn.ReadMessage()
	return data, err
}
