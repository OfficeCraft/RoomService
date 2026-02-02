package websocket

import (
	"log"
	"net/http"

	"github.com/OfficeCraft/RoomService/internal/auth"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func Handler(hub *Hub, w http.ResponseWriter, r *http.Request) {
	log.Println("WebSocket connection attempt")

	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Printf("WebSocket upgrade failed: %s\n", err.Error())
		return
	}

	log.Println("WebSocket connection established")
	cookie, err := r.Cookie("auth_token")
	if err != nil {
		log.Printf("Failed to get auth_token cookie: %s\n", err.Error())
		conn.Close()
		return
	}

	userId, err := auth.ParseJWTTokenToUserID(cookie.Value, "topsecretkey")

	if err != nil {
		log.Printf("Failed to parse JWT token: %s\n", err.Error())
		conn.Close()
		return
	}

	log.Printf("Authenticated user: %s", userId)

	roomId := r.URL.Query().Get("roomId")

	if roomId == "" {
		log.Printf("No roomId provided in query parameters")
		conn.Close()
		return
	}

	if !hub.Rooms.RoomExists(roomId) {
		log.Printf("Room %s does not exist", roomId)
		conn.Close()
		return
	}

	client := &Client{
		Id:     userId,
		Conn:   conn,
		RoomId: roomId,
		Hub:    hub,
		Send:   make(chan []byte, 256),
	}

	hub.Register <- client

	// defer conn.Close()
	go client.WritePump()
	go client.ReadPump()
}
