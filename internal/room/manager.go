package room

import (
	"fmt"
	"sync"

	"github.com/google/uuid"
)

type Manager struct {
	rooms map[string]*Room
	mu    sync.Mutex
}

func NewManager() *Manager {
	return &Manager{
		rooms: make(map[string]*Room),
	}
}

func (m *Manager) CreateRoom() (*Room, string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for {
		id := uuid.New().String()

		if _, exists := m.rooms[id]; !exists {
			room := NewRoom(id)
			m.rooms[id] = room
			return room, id
		}
	}
}

func (m *Manager) DeleteRoom(id string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.rooms, id)
}

func (m *Manager) GetRoom(id string) (*Room, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	room, exists := m.rooms[id]
	return room, exists
}

func (m *Manager) ListRooms() []string {
	m.mu.Lock()
	defer m.mu.Unlock()

	ids := make([]string, 0, len(m.rooms))
	for id := range m.rooms {
		ids = append(ids, id)
	}
	return ids
}

func (m *Manager) PrintAllRoomsForDebug() {
	m.mu.Lock()
	defer m.mu.Unlock()

	fmt.Println("All Rooms:")
	for id, room := range m.rooms {
		fmt.Println("Room ID:", id)
		room.printUnsafe()
	}
}

func (m *Manager) RoomExists(id string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	_, exists := m.rooms[id]
	return exists
}
