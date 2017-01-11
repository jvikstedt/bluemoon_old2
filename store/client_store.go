package store

import (
	"fmt"
	"sync"

	"github.com/jvikstedt/bluemoon/bluemoon"
)

type ClientStore struct {
	users map[int]bluemoon.Client
	uLock sync.RWMutex
}

func NewClientStore() *ClientStore {
	return &ClientStore{
		users: make(map[int]bluemoon.Client),
	}
}

func (us *ClientStore) Add(user bluemoon.Client) error {
	us.uLock.Lock()
	defer us.uLock.Unlock()

	if _, ok := us.users[user.ID()]; ok {
		return fmt.Errorf("users map already has a record with id of %d", user.ID())
	}

	us.users[user.ID()] = user

	return nil
}

func (us *ClientStore) Remove(user bluemoon.Client) error {
	us.uLock.Lock()
	defer us.uLock.Unlock()

	if _, ok := us.users[user.ID()]; !ok {
		return fmt.Errorf("user not found with id of %d", user.ID())
	}

	delete(us.users, user.ID())

	return nil
}

func (us *ClientStore) ByID(id int) (bluemoon.Client, error) {
	us.uLock.RLock()
	defer us.uLock.RUnlock()

	if v, ok := us.users[id]; ok {
		return v, nil
	}

	return nil, fmt.Errorf("User with id of %d not found", id)
}

func (us *ClientStore) One() (bluemoon.Client, error) {
	us.uLock.RLock()
	defer us.uLock.RUnlock()

	for _, user := range us.users {
		return user, nil
	}
	return nil, fmt.Errorf("No clients found")
}
