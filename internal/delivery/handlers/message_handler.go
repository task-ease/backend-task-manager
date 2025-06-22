package handlers

import (
	"github.com/gin-gonic/gin"
	"go-postgres-test/infrastructure/auth"
	"go-postgres-test/internal/middleware"
	"go-postgres-test/internal/usecase"
	"net/http"
)

type MessageHandler struct {
	uc *usecase.MessageUsecase
}

func NewMessageHandler(uc *usecase.MessageUsecase) *MessageHandler {
	return &MessageHandler{uc: uc}
}

func (h *MessageHandler) RegisterRoutes(router *gin.RouterGroup) {
	authService := auth.NewJWTService()

	protected := router.Group("/message", middleware.JWTMiddleware(authService))

	protected.GET("/get-all-messages/:chatId", h.GetAllMessages)
}

func (h *MessageHandler) GetAllMessages(c *gin.Context) {
	chatId := c.Param("chatId")
	messageList, err := h.uc.GetAllMessages(chatId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, messageList)
}
