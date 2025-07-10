package usecase

import (
	"github.com/google/uuid"
	"go-postgres-test/internal/domain"
	"go-postgres-test/internal/response"
	"go-postgres-test/internal/types/user"
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

func (c *ChatUsecase) AddUserToChat(userId uuid.UUID, chatId string, workspaceId uuid.UUID, role user.ChatRole) error {
	return c.chatRepo.AddUserToChat(userId, chatId, workspaceId, role)
}

func (c *ChatUsecase) GetAllUserChats(userId uuid.UUID, workspaceId uuid.UUID) ([]response.GetChats, error) {
	return c.chatRepo.GetAllUserChats(userId, workspaceId)
}

func (c *ChatUsecase) GetChatsBySearch(userID uuid.UUID, workspaceId uuid.UUID, value string) ([]response.GetChatsSearch, error) {
	return c.chatRepo.GetChatsBySearch(userID, workspaceId, value)
}
