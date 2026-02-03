package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/OfficeCraft/RoomService/internal/room"
	"github.com/OfficeCraft/RoomService/internal/websocket"
)

var roomManager = room.NewManager()
var hub = websocket.NewHub(roomManager)

type RoomCreateResponse struct {
	RoomID  string `json:"room_id"`
	Success bool   `json:"success"`
}

func Start(addr string) {

	go hub.Run()

	mux := http.NewServeMux()

	mux.HandleFunc("/ws/room", func(w http.ResponseWriter, r *http.Request) {
		log.Println("WebSocket /ws/room endpoint hit")
		websocket.Handler(hub, w, r)
	})

	mux.HandleFunc("/ping", pingHandler)

	mux.HandleFunc("/room/create", createRoomHandler)
	mux.HandleFunc("/rooms", listRoomsHandler)
	mux.HandleFunc("/room/getClients", getListofClientsInRoomHandler)

	log.Printf("Starting server on %s\n", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("Server failed to start: %s\n", err.Error())
	}
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("pong"))
}

func createRoomHandler(w http.ResponseWriter, r *http.Request) {
	// Placeholder for room creation logic
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	_, id := roomManager.CreateRoom()
	log.Printf("Room created with ID: %s\n", id)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	response := RoomCreateResponse{
		RoomID:  id,
		Success: true,
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Failed to create room", http.StatusInternalServerError)
		return
	}

	w.Write(jsonResponse)
}

func listRoomsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	rooms := roomManager.ListRooms()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	jsonResponse, err := json.Marshal(rooms)
	if err != nil {
		http.Error(w, "Failed to list rooms", http.StatusInternalServerError)
		return
	}

	w.Write(jsonResponse)
}

func getListofClientsInRoomHandler(w http.ResponseWriter, r *http.Request) {
	roomId := r.URL.Query().Get("roomId")
	room, exists := roomManager.GetRoom(roomId)
	if !exists {
		http.Error(w, "Room not found", http.StatusNotFound)
		return
	}

	clients := room.ListClients()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	jsonResponse, err := json.Marshal(clients)
	if err != nil {
		http.Error(w, "Failed to get clients", http.StatusInternalServerError)
		return
	}

	w.Write(jsonResponse)
}
