package socket

import (
	"fmt"
	"net"

	"github.com/jvikstedt/bluemoon/gate"
)

type Server struct {
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Listen(addr string) error {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", addr)
	if err != nil {
		return err
	}
	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		return err
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		cw := NewConnectionWrapper(conn)
		go s.handleSocket(cw)
	}
	return nil
}

func (s *Server) handleSocket(cw *ConnectionWrapper) {
	defer cw.Close()
	worker := gate.NewWorker(1, cw, func(worker *gate.Worker, data []byte) {
		fmt.Printf("message %s\n", string(data))
		worker.Write([]byte("pong"))
	})
	go worker.EnableReader()
	worker.EnableWriter()
}
