package websocket

import (
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	Id     string
	Conn   *websocket.Conn
	RoomId string
	Send   chan []byte
	Hub    *Hub
}

func (c *Client) ReadPump() {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()

	for {
		var payload struct {
			X float64 `json:"x"`
			Y float64 `json:"y"`
		}

		err := c.Conn.ReadJSON(&payload)
		if err != nil {
			log.Println("Error reading JSON:", err)
			break
		}

		c.Hub.Broadcast <- Message{
			RoomId:   c.RoomId,
			ClientId: c,
			X:        payload.X,
			Y:        payload.Y,
		}

	}
}

func (c *Client) WritePump() {
	defer c.Conn.Close()

	for msg := range c.Send {
		err := c.Conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			log.Println("Write error:", err)
			return
		}
	}
}
