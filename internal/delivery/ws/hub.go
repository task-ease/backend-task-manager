package ws

import (
	"github.com/google/uuid"
	"go-postgres-test/internal/types"
)

type Hub struct {
	Clients    map[*Client]bool
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan Message
}

type Message struct {
	RoomID string
	UserID uuid.UUID
	Data   string
	Type   types.WsMessageTypes
}

func NewHub() *Hub {
	return &Hub{
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
			h.Clients[client] = true

		case client := <-h.Unregister:
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				//h.Broadcast <- Message{
				//	RoomID: client.RoomID,
				//	Data:   "disconnected",
				//	UserID: client.ID,
				//	Type:   types.TypeStatus,
				//}
				close(client.Send)
			}

		case message := <-h.Broadcast:
			for client := range h.Clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(h.Clients, client)
				}
			}
		}

	}
}
