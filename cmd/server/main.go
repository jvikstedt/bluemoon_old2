package main

import (
	"log"
	"net"
	"os"

	"github.com/jvikstedt/bluemoon/game"
	"github.com/jvikstedt/bluemoon/msg"
	"github.com/jvikstedt/bluemoon/net/socket"
)

var hub *game.Hub

func main() {
	logger := log.New(os.Stdout, "", log.LstdFlags)

	packer := msg.NewPacker(';')
	hub = game.NewHub(logger, packer)

	server := socket.NewServer(logger, ManageConnFunc)
	if err := server.Listen(":8080"); err != nil {
		panic(err)
	}
}

func ManageConnFunc(conn *net.TCPConn) {
	hub.HandleConnection(conn)
}
