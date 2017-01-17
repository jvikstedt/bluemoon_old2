package room

import (
	"fmt"
	"time"

	"github.com/jvikstedt/bluemoon/bm"
	"github.com/jvikstedt/bluemoon/worker"
)

type Game struct {
	log      bm.Logger
	id       int
	hub      *worker.Hub
	running  bool
	entities map[int]worker.Entity
	eventCh  chan worker.Event
	users    map[int]*worker.User
}

func NewGame(log bm.Logger, id int, hub *worker.Hub) *Game {
	return &Game{
		log:      log,
		id:       id,
		hub:      hub,
		running:  true,
		entities: make(map[int]worker.Entity),
		eventCh:  make(chan worker.Event, 20),
		users:    make(map[int]*worker.User),
	}
}

func (r *Game) ID() int {
	return r.id
}

func (r *Game) Type() string {
	return "game"
}

func (r *Game) Run() {
	tickChan := time.NewTicker(time.Millisecond * 100).C

	last := time.Now()
	for r.running {
		select {
		case <-tickChan:
			delta := time.Since(last).Seconds() + 1
			last = time.Now()

			for _, v := range r.entities {
				changed, err := v.Update(delta)
				if err != nil {
					r.log.Warnln(fmt.Sprintf("Error while executing a entity Update: %s", err.Error()))
				}
				if changed {
					r.hub.Broadcast([]byte(fmt.Sprintf(`{"name": "move", "id": %d, "x": %d, "y": %d}`, v.ID(), v.X(), v.Y())))
				}
			}
		case e := <-r.eventCh:
			err := e.Execute(r)
			if err != nil {
				r.log.Warnln(fmt.Sprintf("Error while executing a event: %s", err.Error()))
			}
		}
	}
}

func (r *Game) Broadcast(data []byte) {
	ids := make([]int, len(r.users))
	i := 0
	for k := range r.users {
		ids[i] = k
		i++
	}

	r.hub.BroadcastTo(ids, data)
}

func (r *Game) BroadcastTo(userIds []int, payload []byte) {
	r.hub.BroadcastTo(userIds, payload)
}

func (r *Game) AddEvent(e worker.Event) {
	r.eventCh <- e
}

func (r *Game) AddEntity(e worker.Entity) {
	r.entities[e.ID()] = e
}

func (r *Game) EntityById(id int) worker.Entity {
	return r.entities[id]
}

func (r *Game) Entities() map[int]worker.Entity {
	return r.entities
}

func (r *Game) RemoveEntity(e worker.Entity) {
	delete(r.entities, e.ID())
}

func (r *Game) RemoveEntityByID(id int) {
	delete(r.entities, id)
}

func (r *Game) AddUser(u *worker.User) {
	r.users[u.ID()] = u
}

func (r *Game) RemoveUser(u *worker.User) {
	delete(r.users, u.ID())
}

func (r *Game) RemoveUserByID(id int) {
	delete(r.users, id)
}