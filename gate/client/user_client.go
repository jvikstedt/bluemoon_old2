package client

import "github.com/jvikstedt/bluemoon/bluemoon"

type HandleUserDataFunc func(user *User, data []byte)

type User struct {
	id         int
	writeCh    chan []byte
	doneCh     chan bool
	rw         bluemoon.ReadWriter
	handleData HandleUserDataFunc
}

func (c *User) EnableWriter() {
	for {
		select {
		case data := <-c.writeCh:
			c.rw.Write(data)
		case <-c.doneCh:
			return
		}
	}
}

func (c *User) EnableReader() {
	for {
		if data, err := c.rw.Read(); err != nil {
			break
		} else {
			c.handleData(c, data)
		}
	}
	c.Close()
}

func (c *User) Write(data []byte) {
	c.writeCh <- data
}

func (c *User) ID() int {
	return c.id
}

func (c *User) Close() {
	c.doneCh <- true
}

func NewUser(id int, rw bluemoon.ReadWriter, dh HandleUserDataFunc) *User {
	return &User{
		id:         id,
		writeCh:    make(chan []byte, 5),
		doneCh:     make(chan bool, 3),
		rw:         rw,
		handleData: dh,
	}
}
