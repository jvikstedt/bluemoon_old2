package client

import "github.com/jvikstedt/bluemoon/bluemoon"

type WorkerClient struct {
	*bluemoon.BaseClient
}

func NewWorkerClient(id int, rw bluemoon.ReadWriter, dh bluemoon.HandleClientDataFunc) *WorkerClient {
	return &WorkerClient{
		BaseClient: bluemoon.NewBaseClient(id, rw, dh),
	}
}
