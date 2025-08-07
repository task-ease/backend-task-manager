package usecase

import (
	"go-postgres-test/internal/domain"
	"go-postgres-test/internal/enums"
	"go-postgres-test/internal/response"

	"github.com/google/uuid"
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

func (c *ChatUsecase) AddUserToChat(userId uuid.UUID, chatId string, workspaceId uuid.UUID, role enums.ChatRole) error {
	return c.chatRepo.AddUserToChat(userId, chatId, workspaceId, role)
}

func (c *ChatUsecase) GetAllUserChats(userId uuid.UUID, workspaceId uuid.UUID) ([]response.GetChats, error) {
	return c.chatRepo.GetAllUserChats(userId, workspaceId)
}

func (c *ChatUsecase) GetChatsBySearch(userID uuid.UUID, workspaceId uuid.UUID, value string) ([]response.GetChatsSearch, error) {
	return c.chatRepo.GetChatsBySearch(userID, workspaceId, value)
}
