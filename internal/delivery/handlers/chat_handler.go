package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go-postgres-test/infrastructure/auth"
	"go-postgres-test/internal/domain"
	"go-postgres-test/internal/middleware"
	"go-postgres-test/internal/usecase"
	"net/http"
)

type ChatHandler struct {
	uc usecase.ChatUsecase
}

func NewChatHandler(uc usecase.ChatUsecase) *ChatHandler {
	return &ChatHandler{uc: uc}
}

func (h *ChatHandler) RegisterRoutes(router *gin.RouterGroup) {
	authService := auth.NewJWTService()

	protected := router.Group("/chat", middleware.JWTMiddleware(authService))

	protected.GET("/get-all-user-chats/:workspaceId", h.GetAllUserChats)
	protected.GET("/get-all-user-chats-search/:value", h.GetChatsBySearch)

	protected.POST("/create-chat/:participantId", h.CreateChat)
	protected.POST("/add-user-to-chat", h.AddUserToChat)
}

func (h *ChatHandler) CreateChat(c *gin.Context) {
	var chat domain.Chat

	if err := c.ShouldBindJSON(&chat); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

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

	pIdStr := c.Param("participantId")
	pId, err := uuid.Parse(pIdStr)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	chat.CreatorID = userId
	if err = h.uc.CreateChat(&chat, pId); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"chat": chat})
}

func (h *ChatHandler) AddUserToChat(c *gin.Context) {
	var input struct {
		ChatID      string    `json:"chatId"`
		UserID      uuid.UUID `json:"userId"`
		WorkspaceID uuid.UUID `json:"workspaceId"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	err := h.uc.AddUserToChat(input.UserID, input.ChatID, input.WorkspaceID)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "success"})
}

func (h *ChatHandler) GetAllUserChats(c *gin.Context) {
	workspaceId, err := uuid.Parse(c.Param("workspaceId"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

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

	chatList, err := h.uc.GetAllUserChats(userId, workspaceId)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, chatList)
}

func (h *ChatHandler) GetChatsBySearch(c *gin.Context) {
	value := c.Param("value")
	workspaceId, err := uuid.Parse(c.Query("workspaceId"))
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

	chatList, err := h.uc.GetChatsBySearch(userId, workspaceId, value)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, chatList)
}
