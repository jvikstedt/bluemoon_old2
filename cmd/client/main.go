package main

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/jvikstedt/bluemoon/bm"
	"github.com/jvikstedt/bluemoon/socket"
	"github.com/jvikstedt/bluemoon/ws"
)

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
	fmt.Printf("Give address\n")
	fmt.Printf("Example websocket: ws://localhost:4000/ws\n")
	fmt.Printf("Example socket: :5000\n")
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

func handleUserInput(writer bm.Writer, quitChan chan bool) {
	reader := bufio.NewReader(os.Stdin)

Loop:
	for {
		select {
		case <-quitChan:
			break Loop
		default:
			fmt.Printf("What you want to send to the server?\n")
			userInput, _ := reader.ReadString('\n')
			writer.Write([]byte(userInput))
		}
	}
}
