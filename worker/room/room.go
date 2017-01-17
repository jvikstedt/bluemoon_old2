package room

import (
	"fmt"
	"time"
)

type Entity interface {
	Update(delta float64) (bool, error)
	ID() int
	X() int
	Y() int
}

type Event interface {
	Execute(room *Room) error
}

type Room struct {
	hub      *Hub
	running  bool
	entities map[int]Entity
	eventCh  chan Event
	users    map[int]*User
}

func NewRoom(hub *Hub) *Room {
	return &Room{
		hub:      hub,
		running:  true,
		entities: make(map[int]Entity),
		eventCh:  make(chan Event, 20),
		users:    make(map[int]*User),
	}
}

func (r *Room) ID() int {
	return 1
}

func (r *Room) Run() {
	tickChan := time.NewTicker(time.Millisecond * 100).C

	last := time.Now()
	for r.running {
		select {
		case <-tickChan:
			delta := time.Since(last).Seconds() + 1
			last = time.Now()

			for _, v := range r.entities {
				changed, _ := v.Update(delta)
				if changed {
					r.hub.Broadcast([]byte(fmt.Sprintf(`{"name": "move", "id": %d, "x": %d, "y": %d}`, v.ID(), v.X(), v.Y())))
				}
			}
		case e := <-r.eventCh:
			err := e.Execute(r)
			if err != nil {
				fmt.Printf("Error while executing event: %s\n", err.Error())
			}
		}
	}
}

func (r *Room) Broadcast(data []byte) {
	ids := make([]int, len(r.users))
	i := 0
	for k := range r.users {
		ids[i] = k
		i++
	}

	r.hub.BroadcastTo(ids, data)
}

func (r *Room) BroadcastTo(userIds []int, payload []byte) {
	r.hub.BroadcastTo(userIds, payload)
}

func (r *Room) AddEvent(e Event) {
	r.eventCh <- e
}

func (r *Room) AddEntity(e Entity) {
	r.entities[e.ID()] = e
}

func (r *Room) EntityById(id int) Entity {
	return r.entities[id]
}

func (r *Room) Entities() map[int]Entity {
	return r.entities
}

func (r *Room) RemoveEntity(e Entity) {
	delete(r.entities, e.ID())
}

func (r *Room) RemoveEntityByID(id int) {
	delete(r.entities, id)
}

func (r *Room) AddUser(u *User) {
	r.users[u.ID()] = u
}

func (r *Room) RemoveUser(u *User) {
	delete(r.users, u.ID())
}

func (r *Room) RemoveUserByID(id int) {
	delete(r.users, id)
}
