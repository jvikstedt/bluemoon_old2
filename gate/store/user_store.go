package store

import (
	"fmt"
	"sync"

	"github.com/jvikstedt/bluemoon/gate/client"
)

type UserStore struct {
	users map[int]*client.User
	uLock sync.RWMutex
}

func NewUserStore() *UserStore {
	return &UserStore{
		users: make(map[int]*client.User),
	}
}

func (us *UserStore) Add(user *client.User) error {
	us.uLock.Lock()
	defer us.uLock.Unlock()

	if _, ok := us.users[user.ID()]; ok {
		return fmt.Errorf("users map already has a record with id of %d", user.ID())
	}

	us.users[user.ID()] = user

	return nil
}

func (us *UserStore) Remove(user *client.User) error {
	us.uLock.Lock()
	defer us.uLock.Unlock()

	if _, ok := us.users[user.ID()]; !ok {
		return fmt.Errorf("user not found with id of %d", user.ID())
	}

	delete(us.users, user.ID())

	return nil
}

func (us *UserStore) ByID(id int) (*client.User, error) {
	us.uLock.RLock()
	defer us.uLock.RUnlock()

	if v, ok := us.users[id]; ok {
		return v, nil
	}

	return nil, fmt.Errorf("User with id of %d not found", id)
}
