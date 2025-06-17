package domain

import (
	"github.com/google/uuid"
	"go-postgres-test/internal/types"
	"time"
)

type Message struct {
	ID          string            `db:"id"`
	ChatID      string            `db:"chat_id"`
	SenderID    uuid.UUID         `db:"sender_id"`
	Content     string            `db:"content"`
	MessageType types.MessageType `db:"message_type"`
	CreatedAt   time.Time         `db:"created_at"`
	UpdatedAt   time.Time         `db:"updated_at"`
	IsEdited    bool              `db:"is_edited"`
	IsDeleted   bool              `db:"is_deleted"`
	ReplyTo     *uuid.UUID        `db:"reply_to"`
}

type MessageRepository interface {
	AddMessage(message *Message) error
	GetAllMessages(chatId string, userId uuid.UUID) ([]*Message, error)
}
