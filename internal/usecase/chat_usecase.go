package usecase

import (
	"github.com/google/uuid"
	"go-postgres-test/internal/domain"
	"go-postgres-test/internal/response"
)

type ChatUsecase struct {
	chatRepo domain.ChatRepository
}

func NewChatUsecase(chatRepo domain.ChatRepository) *ChatUsecase {
	return &ChatUsecase{chatRepo: chatRepo}
}

func (c *ChatUsecase) CreateChat(chat *domain.Chat, participantId uuid.UUID) error {
	return c.chatRepo.CreateChat(chat, participantId)
}

func (c *ChatUsecase) AddUserToChat(userId uuid.UUID, chatId string, workspaceId uuid.UUID) error {
	return c.chatRepo.AddUserToChat(userId, chatId, workspaceId)
}

func (c *ChatUsecase) GetAllUserChats(userId uuid.UUID, workspaceId uuid.UUID) ([]response.GetChats, error) {
	return c.chatRepo.GetAllUserChats(userId, workspaceId)
}
