package client

import "github.com/jvikstedt/bluemoon/bluemoon"

type UserClient struct {
	*bluemoon.BaseClient
}

func NewUserClient(id int, rw bluemoon.ReadWriter, dh bluemoon.HandleClientDataFunc) *UserClient {
	return &UserClient{
		BaseClient: bluemoon.NewBaseClient(id, rw, dh),
	}
}
