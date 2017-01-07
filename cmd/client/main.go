package main

import (
	"fmt"
	"time"

	"github.com/jvikstedt/bluemoon/socket"
	"github.com/jvikstedt/bluemoon/ws"
)

type ReadWriter interface {
	Write(data []byte)
	Read() ([]byte, error)
}

var protocols = map[int]string{
	1: "websocket",
	2: "socket",
}

func main() {
	var protocolID int
	for {
		fmt.Printf("Select protocol:\n")
		for i, val := range protocols {
			fmt.Printf("%d: %s\n", i, val)
		}
		fmt.Scanf("%d", &protocolID)
		if _, ok := protocols[protocolID]; ok {
			break
		}
	}
	fmt.Printf("Selected protocol: %s\n", protocols[protocolID])
	var addr string
	fmt.Printf("Give address (ex: :4000)\n")
	fmt.Scanf("%s", &addr)

	for {
		fmt.Printf("Trying to connect to: %s with protocol: %s\n", addr, protocols[protocolID])

		switch protocolID {
		case 1:
			connectViaWebsocket(addr)
		case 2:
			connectViaSocket(addr)
		}

		time.Sleep(time.Second * 1)
	}
}

func connectViaWebsocket(addr string) {
	wsClient := ws.NewClient()
	conn, err := wsClient.Connect(addr)
	if err != nil {
		fmt.Println(err)
		return
	}

	cw := ws.NewConnectionWrapper(conn)
	defer cw.Close()

	quitChan := make(chan bool)

	go handleUserInput(cw, quitChan)

	for {
		data, err := cw.Read()
		if err != nil {
			fmt.Println(err)
			quitChan <- true
			break
		}
		fmt.Printf("From server: %s\n", data)
	}
}

func connectViaSocket(addr string) {
	sClient := socket.NewClient()
	conn, err := sClient.Connect(addr)
	if err != nil {
		fmt.Println(err)
		return
	}

	cw := socket.NewConnectionWrapper(conn)
	defer cw.Close()

	quitChan := make(chan bool)

	go handleUserInput(cw, quitChan)

	for {
		data, err := cw.Read()
		if err != nil {
			fmt.Println(err)
			quitChan <- true
			break
		}
		fmt.Printf("From server: %s\n", data)
	}
}

func handleUserInput(rw ReadWriter, quitChan chan bool) {
Loop:
	for {
		select {
		case <-quitChan:
			break Loop
		default:
			var userInput string
			fmt.Printf("What you want to send to the server?\n")
			fmt.Scanf("%s", &userInput)
			rw.Write([]byte(userInput))
		}
	}
}
