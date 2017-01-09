package client

import "github.com/jvikstedt/bluemoon/bluemoon"

type HandleWorkerDataFunc func(user *Worker, data []byte)

type Worker struct {
	id         int
	writeCh    chan []byte
	doneCh     chan bool
	rw         bluemoon.ReadWriter
	handleData HandleWorkerDataFunc
}

func (c *Worker) EnableWriter() {
	for {
		select {
		case data := <-c.writeCh:
			c.rw.Write(data)
		case <-c.doneCh:
			return
		}
	}
}

func (c *Worker) EnableReader() {
	for {
		if data, err := c.rw.Read(); err != nil {
			break
		} else {
			c.handleData(c, data)
		}
	}
	c.Close()
}

func (c *Worker) Write(data []byte) {
	c.writeCh <- data
}

func (c *Worker) ID() int {
	return c.id
}

func (c *Worker) Close() {
	c.doneCh <- true
}

func NewWorker(id int, rw bluemoon.ReadWriter, dh HandleWorkerDataFunc) *Worker {
	return &Worker{
		id:         id,
		writeCh:    make(chan []byte, 5),
		doneCh:     make(chan bool, 3),
		rw:         rw,
		handleData: dh,
	}
}
