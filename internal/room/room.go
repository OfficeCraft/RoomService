package room

import (
	"fmt"
	"sync"
)

type PlayerState struct {
	X float64
	Y float64
}

type Room struct {
	id      string
	Players map[string]*PlayerState
	mu      sync.Mutex
}

func NewRoom(id string) *Room {
	return &Room{
		id:      id,
		Players: make(map[string]*PlayerState),
	}
}

func (r *Room) UpdatePlayer(id string, x, y float64) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.Players[id] = &PlayerState{X: x, Y: y}
	r.PrintRoomForDebug()
}

func (r *Room) RemovePlayer(id string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.Players, id)
	r.PrintRoomForDebug()
}

func (r *Room) AddPlayer(id string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.Players[id] = &PlayerState{X: 0, Y: 0}
	r.PrintRoomForDebug()
}

func (r *Room) PrintRoomForDebug() {
	r.mu.Lock()
	defer r.mu.Unlock()
	fmt.Println("Room ID:", r.id)
	for id, state := range r.Players {
		fmt.Println("Player:", id, "X:", state.X, "Y:", state.Y)
	}
}
