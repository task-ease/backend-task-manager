package domain

import (
	"context"
	"go-postgres-test/internal/entities"
	"io"

	"github.com/google/uuid"
)

type MessageRepository interface {
	AddMessage(ctx context.Context, message entities.Message) error
	UploadImage(image io.Reader) (string, error)
	AddAttachment(ctx context.Context, attachment entities.Attachment) error
	SetMessageRead(ctx context.Context, messageId, userId uuid.UUID) error
	GetAllMessages(ctx context.Context, chatId string) ([]entities.Message, error)
	AddMessageReads(ctx context.Context, chatId, messageId string, senderId uuid.UUID) ([]uuid.UUID, error)
	GetAllChatImages(ctx context.Context, chatId string) (*[]entities.Attachment, error)

	AddMessageTx(ctx context.Context, exec entities.Execer, message entities.Message) error
	AddAttachmentTx(ctx context.Context, exec entities.Execer, attachment entities.Attachment) error
	SetMessageReadTx(ctx context.Context, exec entities.Execer, messageId, userId uuid.UUID) error
	AddMessageReadsTx(ctx context.Context, exec entities.Execer, chatId, messageId string, senderId uuid.UUID) ([]uuid.UUID, error)
}
