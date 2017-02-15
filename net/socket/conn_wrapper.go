package socket

import (
	"bufio"
	"net"
)

type ConnWrapper struct {
	conn   *net.TCPConn
	reader *bufio.Reader
	writer *bufio.Writer
}

func NewConnWrapper(conn *net.TCPConn) *ConnWrapper {
	return &ConnWrapper{
		conn:   conn,
		reader: bufio.NewReader(conn),
		writer: bufio.NewWriter(conn),
	}
}

func (c *ConnWrapper) Read(p []byte) (n int, err error) {
	p, err = c.reader.ReadBytes('\n')

	if p[len(p)-1] == '\n' {
		p = p[:len(p)-1]
	}

	return len(p), err
}

func (c *ConnWrapper) Write(p []byte) (n int, err error) {
	if p[len(p)-1] != '\n' {
		p = append(p, '\n')
	}
	n, err = c.writer.Write(p)
	if err != nil {
		return
	}
	return n, c.writer.Flush()
}
