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
	actions  []worker.Action
	actionCh chan worker.Action
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
		actions:  []worker.Action{},
		actionCh: make(chan worker.Action, 20),
	}
}

func (r *Game) ID() int {
	return r.id
}

func (r *Game) Type() string {
	return "game"
}

func (r *Game) Run() {
	tickChan := time.NewTicker(time.Millisecond * 33).C

	last := time.Now()
	for r.running {
		select {
		case <-tickChan:
			delta := time.Since(last).Seconds()
			last = time.Now()

			for _, v := range r.actions {
				v.Run(r, delta)
			}

			for _, v := range r.entities {
				changed, err := v.Update(delta * 10)
				if err != nil {
					r.log.Warnln(fmt.Sprintf("Error while executing a entity Update: %s", err.Error()))
				}
				if changed {
					r.Broadcast(struct {
						Name string `json:"name"`
						ID   int    `json:"id"`
						X    int    `json:"x"`
						Y    int    `json:"y"`
					}{"move", v.ID(), v.X(), v.Y()})
				}
			}
		case e := <-r.eventCh:
			err := e.Execute(r)
			if err != nil {
				r.log.Warnln(fmt.Sprintf("Error while executing a event: %s", err.Error()))
			}
		case a := <-r.actionCh:
			r.actions = append(r.actions, a)
		}
	}
}

func (r *Game) Broadcast(data interface{}) {
	ids := make([]int, len(r.users))
	i := 0
	for k := range r.users {
		ids[i] = k
		i++
	}

	r.hub.BroadcastTo(ids, data)
}

func (r *Game) BroadcastTo(userIds []int, payload interface{}) {
	r.hub.BroadcastTo(userIds, payload)
}

func (r *Game) AddEvent(e worker.Event) {
	r.eventCh <- e
}

func (r *Game) AddAction(a worker.Action) {
	r.actionCh <- a
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
