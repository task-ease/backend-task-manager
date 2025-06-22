package ws

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go-postgres-test/infrastructure/auth"
	"go-postgres-test/internal/domain"
	"go-postgres-test/internal/middleware"
	"go-postgres-test/internal/types"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type WebSocketHandler struct {
	msgRepo domain.MessageRepository
}

type Client struct {
	Conn *websocket.Conn
	Room string
	ID   uuid.UUID
}

type WebSocketMessage struct {
	Type   types.WsMessageTypes `json:"type"`
	UserID uuid.UUID            `json:"userId"`
	Data   string               `json:"data"`
	RoomID string               `json:"roomId"`
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var (
	rooms   = make(map[string][]*Client)
	roomsMu sync.Mutex
)

func NewWebSocketHandler(msgRepo domain.MessageRepository) *WebSocketHandler {
	return &WebSocketHandler{msgRepo: msgRepo}
}

func (h *WebSocketHandler) RegisterRoutes(r *gin.Engine) {
	authService := auth.NewJWTService()

	protected := r.Group("/ws", middleware.JWTMiddleware(authService))

	protected.GET("/:roomId", h.HandleWS)
}

func (h *WebSocketHandler) HandleWS(c *gin.Context) {
	roomId := c.Param("roomId")

	userIdStr, exists := c.Get("userId")

	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user not found"})
		return
	}

	userId, err := uuid.Parse(userIdStr.(string))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()

	client := &Client{Conn: conn, Room: roomId, ID: userId}

	roomsMu.Lock()
	rooms[roomId] = append(rooms[roomId], client)
	roomsMu.Unlock()

	msg := WebSocketMessage{
		Type:   types.TypeConnected,
		UserID: client.ID,
		Data:   "",
		RoomID: roomId,
	}
	sendJSONToRoom(roomId, msg)

	for {
		_, raw, err := conn.ReadMessage()
		if err != nil {
			msg := WebSocketMessage{
				Type:   types.TypeDisconnected,
				UserID: client.ID,
				Data:   "",
				RoomID: roomId,
			}
			sendJSONToRoom(roomId, msg)
			break
		}

		var wsMsg domain.Message
		_ = json.Unmarshal(raw, &wsMsg)

		addData := domain.Message{
			ID:          "",
			ChatID:      roomId,
			SenderID:    client.ID,
			Content:     wsMsg.Content,
			MessageType: wsMsg.MessageType,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		_ = h.msgRepo.AddMessage(&addData)

		dataStr, _ := json.Marshal(addData)

		msg := WebSocketMessage{
			Type:   types.TypeMessage,
			UserID: client.ID,
			Data:   string(dataStr),
			RoomID: roomId,
		}
		sendJSONToRoom(roomId, msg)
	}

	roomsMu.Lock()
	clients := rooms[roomId]
	for i, c := range clients {
		if c == client {
			rooms[roomId] = append(clients[:i], clients[i+1:]...)
			break
		}
	}
	roomsMu.Unlock()
}

func sendJSONToRoom(roomId string, msg WebSocketMessage) {
	roomsMu.Lock()
	defer roomsMu.Unlock()

	for _, c := range rooms[roomId] {
		err := c.Conn.WriteJSON(msg)
		if err != nil {
			log.Println("Error: ", err)
		}
	}
}
