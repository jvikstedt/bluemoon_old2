package main

import (
	"fmt"
	"net"
	"time"

	"github.com/jvikstedt/bluemoon/gate"
)

func main() {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", ":5000")
	if err != nil {
		panic(err)
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		panic(err)
	}
	s := newSocket(conn)
	defer s.Close()

	worker := gate.NewWorker(1, s, func(worker *gate.Worker, data []byte) {
		fmt.Printf("message %s\n", string(data))
	})
	go worker.EnableReader()
	go worker.EnableWriter()
	for {
		worker.Write([]byte("ping"))
		time.Sleep(time.Second * 1)
	}
}

type socket struct {
	conn net.Conn
}

func newSocket(conn net.Conn) *socket {
	return &socket{
		conn: conn,
	}
}

func (s *socket) Close() error {
	return s.conn.Close()
}

func (s *socket) Write(data []byte) {
	s.conn.Write(data)
}

func (s *socket) Read() ([]byte, error) {
	request := make([]byte, 128)
	_, err := s.conn.Read(request)
	return request, err
}
