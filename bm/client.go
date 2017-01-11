package bm

type BaseClient struct {
	id         int
	writeCh    chan []byte
	doneCh     chan bool
	rw         ReadWriter
	handleData HandleClientDataFunc
}

func (c *BaseClient) EnableWriter() {
	for {
		select {
		case data := <-c.writeCh:
			c.rw.Write(data)
		case <-c.doneCh:
			return
		}
	}
}

func (c *BaseClient) EnableReader() {
	for {
		if data, err := c.rw.Read(); err != nil {
			break
		} else {
			c.handleData(c, data)
		}
	}
	c.Close()
}

func (c *BaseClient) Write(data []byte) {
	c.writeCh <- data
}

func (c *BaseClient) ID() int {
	return c.id
}

func (c *BaseClient) Close() {
	c.doneCh <- true
}

func NewBaseClient(id int, rw ReadWriter, dh HandleClientDataFunc) *BaseClient {
	return &BaseClient{
		id:         id,
		writeCh:    make(chan []byte, 5),
		doneCh:     make(chan bool, 3),
		rw:         rw,
		handleData: dh,
	}
}
