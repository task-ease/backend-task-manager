package handlers

import (
	"go-postgres-test/infrastructure/auth"
	"go-postgres-test/internal/domain"
	"go-postgres-test/internal/enums"
	"go-postgres-test/internal/middleware"
	"go-postgres-test/internal/request"
	"go-postgres-test/internal/usecase"
	"go-postgres-test/mixins"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ChatHandler struct {
	uc *usecase.ChatUsecase
}

func NewChatHandler(uc *usecase.ChatUsecase) *ChatHandler {
	return &ChatHandler{uc: uc}
}

func (h *ChatHandler) RegisterRoutes(router *gin.RouterGroup) {
	authService := auth.NewJWTService()

	protected := router.Group("/chat", middleware.JWTMiddleware(authService))

	protected.GET("/get-all-user-chats/:workspaceId", h.getAllUserChats)
	protected.GET("/get-all-user-chats-search/:value", h.getChatsBySearch)

	protected.POST("/create-chat/:participantId", h.createChat)
	protected.POST("/add-user-to-chat", h.addUserToChat)
}

func (h *ChatHandler) createChat(c *gin.Context) {
	var chat domain.Chat
	var err error

	if err = c.ShouldBindJSON(&chat); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	chat.CreatorID, err = mixins.ParseUserId(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	participantId, err := mixins.ParamToUUID(c, "participantId")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err = h.uc.CreateChat(c.Request.Context(), &chat, participantId); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"chat": chat})
}

func (h *ChatHandler) addUserToChat(c *gin.Context) {
	var input request.AddUserToChat

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	if err := h.uc.AddUserToChat(c.Request.Context(), input.UserID, input.ChatID, input.WorkspaceID, enums.ChatUser); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func (h *ChatHandler) getAllUserChats(c *gin.Context) {
	workspaceId, err := mixins.ParamToUUID(c, "workspaceId")

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	userId, err := mixins.ParseUserId(c)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	chatList, err := h.uc.GetAllUserChats(c.Request.Context(), userId, workspaceId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, chatList)
}

func (h *ChatHandler) getChatsBySearch(c *gin.Context) {
	value := c.Param("value")
	workspaceId, err := mixins.ParamToUUID(c, "workspaceId")

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userId, err := mixins.ParseUserId(c)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	chatList, err := h.uc.GetChatsBySearch(c.Request.Context(), userId, workspaceId, value)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, chatList)
}
