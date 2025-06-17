package http

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"go-postgres-test/infrastructure/auth"
	"go-postgres-test/internal/delivery/ws"
	"go-postgres-test/internal/middleware"
	"go-postgres-test/internal/types"
	"net/http"
)

type WsHandler struct {
	hub *ws.Hub
}

func NewWsHandler(hub *ws.Hub) *WsHandler {
	return &WsHandler{hub: hub}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *WsHandler) RegisterRoutes(router *gin.Engine) {
	authService := auth.NewJWTService()

	protected := router.Group("/ws", middleware.JWTMiddleware(authService))

	protected.GET("/:roomId", h.ServeWs())
}

func (h *WsHandler) ServeWs() gin.HandlerFunc {
	return func(c *gin.Context) {
		roomID := c.Param("roomId")

		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			return
		}

		userIdStr, exists := c.Get("userId")

		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user id not found"})
			return
		}

		userID, err := uuid.Parse(userIdStr.(string))

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		client := &ws.Client{
			Conn:   conn,
			Send:   make(chan ws.Message, 256),
			RoomID: roomID,
			Hub:    h.hub,
			ID:     userID,
		}

		h.hub.Register <- client

		go client.ReadPump()
		go client.WritePump()

		h.hub.Broadcast <- ws.Message{
			RoomID: roomID,
			UserID: userID,
			Data:   "joined",
			Type:   types.TypeStatus,
		}
	}
}
