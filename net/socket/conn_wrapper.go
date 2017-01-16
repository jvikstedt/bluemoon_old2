package socket

import (
	"bufio"
	"net"
)

type ConnectionWrapper struct {
	conn   net.Conn
	reader *bufio.Reader
	writer *bufio.Writer
}

func NewConnectionWrapper(conn net.Conn) *ConnectionWrapper {
	return &ConnectionWrapper{
		conn:   conn,
		writer: bufio.NewWriter(conn),
		reader: bufio.NewReader(conn),
	}
}

func (cw *ConnectionWrapper) Close() error {
	return cw.conn.Close()
}

func (cw *ConnectionWrapper) Write(data []byte) {
	cw.writer.WriteString(string(data))
	cw.writer.Flush()
}

func (cw *ConnectionWrapper) Read() ([]byte, error) {
	line, _ := cw.reader.ReadString('\n')
	return []byte(line), nil
}
