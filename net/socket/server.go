package socket

import (
	"github.com/jvikstedt/bluemoon/log"
	"net"
)

type ManageConnFunc func(*net.TCPConn)

type Server struct {
	log        log.Logger
	manageConn ManageConnFunc
	running    bool
}

func NewServer(log log.Logger, manageConn ManageConnFunc) *Server {
	return &Server{
		log:        log,
		manageConn: manageConn,
		running:    true,
	}
}

// Listen is a blocking operation that listens for tcp4 connections
// It spawns a goroutine for every connection and calls ManageConnFunc with it
// It handles closing connection
func (s *Server) Listen(addr string) error {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", addr)
	if err != nil {
		return err
	}
	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		return err
	}

	for s.running {
		conn, err := listener.AcceptTCP()
		if err != nil {
			s.log.Printf("Error when accepting client: %s\n", err.Error())
			continue
		}
		go s.handleconn(conn)
	}
	return nil
}

func (s *Server) handleconn(conn *net.TCPConn) {
	defer conn.Close()
	s.manageConn(conn)
}

func (s *Server) Shutdown() {
	s.running = false
}
