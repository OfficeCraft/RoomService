package server

import (
	"log"
	"net/http"

	"github.com/OfficeCraft/RoomService/internal/websocket"
)

func Start(addr string) {
	http.HandleFunc("/ping", pongHandler)

	http.HandleFunc("/ws", websocket.Handler)

	log.Printf("Starting server at %s\n", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Could not start server: %s\n", err.Error())
	}
}

func pongHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("pong"))
}
