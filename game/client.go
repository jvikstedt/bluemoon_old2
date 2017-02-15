package game

import (
	"bufio"

	"github.com/jvikstedt/bluemoon/log"
)

type HandleFunc func(c *Client, bytes []byte)

type Client struct {
	id      int
	log     log.Logger
	conn    Connection
	reader  *bufio.Reader
	writer  *bufio.Writer
	handle  HandleFunc
	open    bool
	writeCh chan []byte
}

func NewClient(log log.Logger, id int, conn Connection, handle HandleFunc) *Client {
	return &Client{
		log:     log,
		id:      id,
		conn:    conn,
		reader:  bufio.NewReader(conn),
		writer:  bufio.NewWriter(conn),
		handle:  handle,
		open:    true,
		writeCh: make(chan []byte, 10),
	}
}

func (c *Client) EnableReader() {
	for c.open {
		bytes, err := c.reader.ReadBytes('\n')
		if err != nil {
			c.log.Printf("Client %d received an error while reading: %v", c.id, err)
			break
		}

		if bytes[len(bytes)-1] == '\n' {
			bytes = bytes[:len(bytes)-1]
		}

		c.handle(c, bytes)
	}
	c.log.Printf("Client %d closing reader\n", c.id)
	c.Close()
}

func (c *Client) EnableWriter() {
	for c.open {
		select {
		case bytes := <-c.writeCh:

			if bytes[len(bytes)-1] != '\n' {
				bytes = append(bytes, '\n')
			}
			_, err := c.writer.Write(bytes)
			if err != nil {
				c.log.Printf("Client %d received an error while writing: %s", c.id, err.Error())
				continue
			}
			err = c.writer.Flush()
			if err != nil {
				c.log.Printf("Client %d received an error while flushing: %s", c.id, err.Error())
			}
		default:
		}
	}
	c.log.Printf("Client %d closing writer\n", c.id)
	c.Close()
}

func (c *Client) Write(bytes []byte) {
	c.writeCh <- bytes
}

func (c *Client) Close() {
	c.open = false
}

func (c *Client) ID() int {
	return c.id
}
