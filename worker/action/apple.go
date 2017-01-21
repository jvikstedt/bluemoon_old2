package action

import (
	"fmt"
	"math/rand"

	"github.com/jvikstedt/bluemoon/worker"
	"github.com/jvikstedt/bluemoon/worker/entity"
	"github.com/jvikstedt/bluemoon/worker/room"
)

type AppleAction struct {
	sinceLastSpawn float64
	index          int
}

func (aa *AppleAction) Run(r worker.Room, delta float64) {
	aa.sinceLastSpawn += delta
	if aa.sinceLastSpawn > 5 || aa.index == 0 {
		aa.sinceLastSpawn = 0
		apple := entity.NewAppleEntity(aa.index, rand.Intn(800), rand.Intn(600))
		aa.index++

		game, ok := r.(*room.Game)
		if !ok {
			fmt.Println("Wrong type of room, expected Game")
			return
		}

		game.AddEntity(apple)

		game.Broadcast(struct {
			Name string `json:"name"`
			ID   int    `json:"id"`
			X    int    `json:"x"`
			Y    int    `json:"y"`
		}{"new_apple", apple.ID(), apple.X(), apple.Y()})
	}
}
