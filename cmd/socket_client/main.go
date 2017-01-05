package main

import (
	"fmt"
	"time"

	"github.com/jvikstedt/bluemoon/gate"
	"github.com/jvikstedt/bluemoon/socket"
)

func main() {
	sClient := socket.NewClient()
	conn, err := sClient.Connect(":5000")
	if err != nil {
		panic(err)
	}

	cw := socket.NewConnectionWrapper(conn)
	defer cw.Close()

	w := gate.NewWorker(1, cw, func(worker *gate.Worker, data []byte) {
		fmt.Printf("Got message: %s\n", string(data))
	})
	go w.EnableReader()
	go w.EnableWriter()

	for {
		w.Write([]byte("Ping"))
		time.Sleep(time.Second * 1)
	}
}
