package usecase

import (
	"github.com/google/uuid"
	"go-postgres-test/internal/domain"
)

type MessageUsecase struct {
	repo domain.MessageRepository
}

func NewMessageUsecase(repo domain.MessageRepository) *MessageUsecase {
	return &MessageUsecase{repo: repo}
}

func (uc *MessageUsecase) AddMessage(message *domain.Message) error {
	return uc.repo.AddMessage(message)
}

func (uc *MessageUsecase) GetAllMessages(chatId string, userId uuid.UUID) ([]*domain.Message, error) {
	return uc.repo.GetAllMessages(chatId, userId)
}
