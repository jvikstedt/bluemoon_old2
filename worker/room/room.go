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
}

func NewRoom(hub *Hub) *Room {
	return &Room{
		hub:      hub,
		running:  true,
		entities: make(map[int]Entity),
		eventCh:  make(chan Event, 20),
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

func (r *Room) Hub() *Hub {
	return r.hub
}

func (r *Room) RemoveEntity(e Entity) {
	delete(r.entities, e.ID())
}
