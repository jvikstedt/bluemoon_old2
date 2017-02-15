package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/jvikstedt/bluemoon/net/socket"
)

var protocols = map[int]string{
	1: "socket",
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
		case 2:
			connectViaWebsocket(addr)
		case 1:
			connectViaSocket(addr)
		}

		time.Sleep(time.Second * 1)
	}
}

func connectViaWebsocket(addr string) {
	panic("not implemented")
}

func connectViaSocket(addr string) {
	tcp, err := socket.Connect(addr)
	if err != nil {
		fmt.Println(err)
		return
	}
	conn := socket.NewConnWrapper(tcp)

	quitChan := make(chan bool)

	go handleUserInput(conn, quitChan)

	var data []byte
	for {
		_, err := conn.Read(data)
		if err != nil {
			fmt.Println(err)
			quitChan <- true
			break
		}
		fmt.Printf("From server: %s\n", string(data))
	}
}

func handleUserInput(writer io.Writer, quitChan chan bool) {
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
