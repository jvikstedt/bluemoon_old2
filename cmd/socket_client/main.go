package main

import (
	"github.com/jvikstedt/bluemoon/socket"
)

func main() {
	sClient := socket.NewClient()
	sClient.Connect(":5000")
}
