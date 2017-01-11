package bm

type Reader interface {
	Read() ([]byte, error)
}

type Writer interface {
	Write(data []byte)
}

type ReadWriter interface {
	Reader
	Writer
}

type Client interface {
	Writer
	ID() int
	Close()
}

type HandleClientDataFunc func(client Client, data []byte)
