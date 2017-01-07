package ws

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

type ManageConnFunc func(*websocket.Conn) error

type Server struct {
	manageConn ManageConnFunc
}

func NewServer(manageConn ManageConnFunc) *Server {
	return &Server{
		manageConn: manageConn,
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err.Error())
		return
	}
	defer conn.Close()
	s.manageConn(conn)
}
