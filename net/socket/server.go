package socket

import "net"

type ManageConnFunc func(*net.TCPConn) error

type Server struct {
	manageConn ManageConnFunc
}

func NewServer(manageConn ManageConnFunc) *Server {
	return &Server{
		manageConn: manageConn,
	}
}

func (s *Server) Listen(addr string) {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", addr)
	if err != nil {
		panic(err)
	}
	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		panic(err)
	}

	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			continue
		}
		s.handleConn(conn)
	}
}

func (s *Server) handleConn(conn *net.TCPConn) error {
	defer conn.Close()
	return s.manageConn(conn)
}
