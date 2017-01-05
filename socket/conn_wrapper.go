package socket

import "net"

type ConnectionWrapper struct {
	conn net.Conn
}

func NewConnectionWrapper(conn net.Conn) *ConnectionWrapper {
	return &ConnectionWrapper{
		conn: conn,
	}
}

func (cw *ConnectionWrapper) Close() error {
	return cw.conn.Close()
}

func (cw *ConnectionWrapper) Write(data []byte) {
	cw.conn.Write(data)
}

func (cw *ConnectionWrapper) Read() ([]byte, error) {
	request := make([]byte, 1024)
	_, err := cw.conn.Read(request)
	return request, err
}
