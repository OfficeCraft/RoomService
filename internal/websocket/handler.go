package websocket

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func Handler(w http.ResponseWriter, r *http.Request) {
	log.Println("WebSocket connection attempt")

	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Printf("WebSocket upgrade failed: %s\n", err.Error())
		return
	}

	log.Println("WebSocket connection established")

	defer conn.Close()

	err = conn.WriteMessage(
		websocket.TextMessage,
		[]byte("Hello connection"),
	)

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Read error: %s\n", err.Error())
			break
		}

		log.Printf("Received: %s\n", message)

		err = conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Printf("Write error: %s\n", err.Error())
			break
		}

	}
}
