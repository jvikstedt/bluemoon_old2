package main

import (
	"fmt"
	"net"

	"github.com/jvikstedt/bluemoon/gate"
)

func listenSockets(addr string) {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", addr)
	if err != nil {
		panic(err)
	}
	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		panic(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		s := newSocket(conn)
		go handleSocket(s)
	}
}

func handleSocket(s *socket) {
	defer s.Close()
	worker := gate.NewWorker(1, s, func(worker *gate.Worker, data []byte) {
		fmt.Printf("message %s\n", string(data))
		worker.Write([]byte("pong"))
	})
	go worker.EnableReader()
	worker.EnableWriter()
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
