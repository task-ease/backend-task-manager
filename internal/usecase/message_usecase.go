package usecase

import (
	"go-postgres-test/internal/domain"
	"io"
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

func (uc *MessageUsecase) GetAllMessages(chatId string) ([]*domain.Message, error) {
	return uc.repo.GetAllMessages(chatId)
}

func (uc *MessageUsecase) UploadImage(image io.Reader) (string, error) {
	return uc.repo.UploadImage(image)
}

func (uc *MessageUsecase) AddAttachment(attachment *domain.Attachment) error {
	return uc.repo.AddAttachment(attachment)
}
