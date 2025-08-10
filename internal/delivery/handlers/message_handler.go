package handlers

import (
	"go-postgres-test/infrastructure/auth"
	"go-postgres-test/internal/dto"
	"go-postgres-test/internal/entities"
	"go-postgres-test/internal/middleware"
	"go-postgres-test/internal/usecase"
	"go-postgres-test/mixins"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

	protected.GET("/get-all-messages/:chatId", h.getAllMessages)
	protected.GET("/chat-images/:chatId", h.getAllChatImages)

	protected.PATCH("/upload-image-list/:chatId", h.uploadImageList)
	protected.PATCH("/add-message", h.addMessage)
	protected.PATCH("/set-read", h.SetMessagesRead)
}

func (h *MessageHandler) getAllMessages(c *gin.Context) {
	messageList, err := h.uc.GetAllMessages(c.Request.Context(), c.Param("chatId"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, messageList)
}

func (h *MessageHandler) uploadImageList(c *gin.Context) {
	var uploadImageDto dto.UploadImage
	var err error

	uploadImageDto.ChatId = c.Param("chatId")
	uploadImageDto.Content = c.Query("ct")

	uploadImageDto.UserId, err = mixins.ParseUserId(c)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid form"})
		return
	}

	uploadImageDto.Form, err = c.MultipartForm()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "user not found"})
		return
	}

	message, err := h.uc.UploadImage(c.Request.Context(), uploadImageDto)

	c.JSON(http.StatusOK, gin.H{"message": message})
}

func (h *MessageHandler) addMessage(c *gin.Context) {
	var message entities.Message
	if err := c.ShouldBindJSON(&message); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	unreadUserIds, err := h.uc.AddMessage(c.Request.Context(), message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"messageInfo": gin.H{
			"messageId":      message.ID,
			"unreadUsersIds": unreadUserIds,
		},
	})
}

func (h *MessageHandler) SetMessagesRead(c *gin.Context) {
	var messages []uuid.UUID
	if err := c.ShouldBindJSON(&messages); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userId, err := mixins.ParseUserId(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.uc.SetMessagesRead(c.Request.Context(), messages, userId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func (h *MessageHandler) getAllChatImages(c *gin.Context) {
	chatId := c.Param("chatId")
	images, err := h.uc.GetAllChatImages(c.Request.Context(), chatId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, images)
}
