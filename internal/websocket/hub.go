package websocket

import (
	"encoding/json"
	"log"

	"github.com/OfficeCraft/RoomService/internal/room"
)

type Message struct {
	RoomId   string
	ClientId *Client
	X        float64
	Y        float64
}

type Hub struct {
	Rooms      *room.Manager
	Clients    map[*Client]bool
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan Message
}

func NewHub(rm *room.Manager) *Hub {
	return &Hub{
		Rooms:      rm,
		Clients:    make(map[*Client]bool),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan Message),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			log.Printf("Registering client %s to room %s", client.Id, client.RoomId)
			h.Clients[client] = true
		case client := <-h.Unregister:
			log.Printf("Unregistering client %s from room %s", client.Id, client.RoomId)
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				close(client.Send)
			}
		case message := <-h.Broadcast:
			log.Printf("Broadcasting message from client %s in room %s", message.ClientId.Id, message.RoomId)
			room, exists := h.Rooms.GetRoom(message.RoomId)
			if !exists {
				log.Printf("Room %s does not exist for broadcasting", message.RoomId)
				continue
			}

			room.UpdatePlayer(message.ClientId.Id, message.X, message.Y)

			update := struct {
				PlayerID string  `json:"player_id"`
				X        float64 `json:"x"`
				Y        float64 `json:"y"`
			}{
				PlayerID: message.ClientId.Id,
				X:        message.X,
				Y:        message.Y,
			}

			data, err := json.Marshal(update)
			if err != nil {
				log.Printf("Failed to marshal update message: %s", err.Error())
				continue
			}

			for client := range h.Clients {
				if client.RoomId == message.RoomId {
					select {
					case client.Send <- data:
						println("Sent message to client", client.Id)
					default:
						close(client.Send)
						delete(h.Clients, client)
					}
				}
			}
		}
	}
}
