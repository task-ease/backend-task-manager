package ws

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Client struct {
	ID     uuid.UUID
	Conn   *websocket.Conn
	Send   chan Message
	RoomID string
	Hub    *Hub
}

func (c *Client) ReadPump() {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()

	for {
		_, msg, err := c.Conn.ReadMessage()
		if err != nil {
			break
		}

		c.Hub.Broadcast <- Message{
			RoomID: c.RoomID,
			Data:   string(msg),
			UserID: c.ID,
		}
	}
}

func (c *Client) WritePump() {
	defer c.Conn.Close()
	for msg := range c.Send {
		msgBytes, _ := json.Marshal(msg)
		if err := c.Conn.WriteMessage(websocket.TextMessage, msgBytes); err != nil {
			break
		}
	}
}
