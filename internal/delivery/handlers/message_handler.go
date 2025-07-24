package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go-postgres-test/infrastructure/auth"
	"go-postgres-test/internal/domain"
	"go-postgres-test/internal/enums"
	"go-postgres-test/internal/middleware"
	"go-postgres-test/internal/usecase"
	"net/http"
	"time"
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
	protected.GET("/chat-images/:chatId", h.GetAllChatImages)
	protected.PATCH("/upload-image-list/:chatId", h.UploadImageList)
	protected.PATCH("/add-message", h.AddMessage)
	protected.PATCH("/set-read", h.SetMessagesRead)
}

func (h *MessageHandler) GetAllMessages(c *gin.Context) {
	chatId := c.Param("chatId")
	messageList, err := h.uc.GetAllMessages(chatId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, messageList)
}

func (h *MessageHandler) UploadImageList(c *gin.Context) {
	chatId := c.Param("chatId")
	content := c.Query("ct")

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

	message := domain.Message{
		ID:          "",
		ChatID:      chatId,
		SenderID:    userId,
		Content:     content,
		MessageType: enums.MessageImage,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}

	unreadUserIds, err := h.uc.AddMessage(&message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to add message"})
		return
	}

	message.UnreadUsersIds = *unreadUserIds

	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid form"})
		return
	}

	files := form.File["images"]

	var messageAttachments []domain.Attachment
	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			continue
		}
		defer file.Close()

		url, err := h.uc.UploadImage(file)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}

		var newAttachment = domain.Attachment{
			ID:         uuid.New(),
			MessageID:  message.ID,
			FileUrl:    url,
			FileType:   enums.MessageImage,
			FileName:   fileHeader.Filename,
			FileSize:   fileHeader.Size,
			UploadedAt: time.Now().UTC(),
			ChatID:     message.ChatID,
		}

		if err := h.uc.AddAttachment(&newAttachment); err != nil {
			continue
		}

		messageAttachments = append(messageAttachments, newAttachment)
	}

	message.Attachments = messageAttachments

	c.JSON(http.StatusOK, gin.H{"message": message})
}

func (h *MessageHandler) AddMessage(c *gin.Context) {
	var message domain.Message
	if err := c.ShouldBindJSON(&message); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	unreadUserIds, err := h.uc.AddMessage(&message)
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
	var messages []domain.Message
	if err := c.ShouldBindJSON(&messages); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userIdStr, exists := c.Get("userId")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user id not found"})
		return
	}

	userId, err := uuid.Parse(userIdStr.(string))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.uc.SetMessagesRead(&messages, userId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, "success")
}

func (h *MessageHandler) GetAllChatImages(c *gin.Context) {
	chatId := c.Param("chatId")
	images, err := h.uc.GetAllChatImages(chatId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, images)
}
