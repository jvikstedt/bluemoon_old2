package bluemoon

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
