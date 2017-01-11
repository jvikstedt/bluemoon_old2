package main

import (
	"fmt"

	"github.com/jvikstedt/bluemoon/bm"
	"github.com/jvikstedt/bluemoon/socket"
)

func main() {
	sClient := socket.NewClient()
	conn, err := sClient.Connect("gate:5000")
	if err != nil {
		panic(err)
	}

	cw := socket.NewConnectionWrapper(conn)
	defer cw.Close()

	w := bm.NewBaseClient(1, cw, func(client bm.Client, data []byte) {
		fmt.Printf("New message from gate: %d\n", client.ID())
		fmt.Print(string(data))
	})

	go w.EnableReader()
	w.EnableWriter()
}
