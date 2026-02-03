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
	Clients    map[string][]*Client
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan Message
}

func NewHub(rm *room.Manager) *Hub {
	return &Hub{
		Rooms:      rm,
		Clients:    make(map[string][]*Client),
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
			h.Clients[client.RoomId] = append(h.Clients[client.RoomId], client)
		case client := <-h.Unregister:
			log.Printf("Unregistering client %s from room %s", client.Id, client.RoomId)
			if clients, ok := h.Clients[client.RoomId]; ok {
				for i, c := range clients {
					if c == client {
						h.Clients[client.RoomId] = append(clients[:i], clients[i+1:]...)
						break
					}
				}
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

			clients := h.Clients[message.RoomId]
			failedClients := []*Client{}

			for _, client := range clients {
				select {
				case client.Send <- data:
				default:
					failedClients = append(failedClients, client)
				}
			}

			// Remove failed clients after iteration
			for _, failedClient := range failedClients {
				close(failedClient.Send)
				for i, c := range h.Clients[message.RoomId] {
					if c == failedClient {
						h.Clients[message.RoomId] = append(h.Clients[message.RoomId][:i], h.Clients[message.RoomId][i+1:]...)
						break
					}
				}
			}
		}
	}
}
