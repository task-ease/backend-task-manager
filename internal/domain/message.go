package domain

import (
	"github.com/google/uuid"
	"go-postgres-test/internal/types"
	"time"
)

type Message struct {
	ID          string            `db:"ID" json:"ID"`
	ChatID      string            `db:"chatID" json:"chatID"`
	SenderID    uuid.UUID         `db:"senderID" json:"senderID"`
	Content     string            `db:"content" json:"content"`
	MessageType types.MessageType `db:"messageType" json:"messageType"`
	CreatedAt   time.Time         `db:"createdAt" json:"createdAt"`
	UpdatedAt   time.Time         `db:"updatedAt" json:"updatedAt"`
	IsEdited    bool              `db:"isEdited" json:"isEdited"`
	IsDeleted   bool              `db:"isDeleted" json:"isDeleted"`
	ReplyTo     *uuid.UUID        `db:"replyTo" json:"replyTo"`
}

type MessageRepository interface {
	AddMessage(message *Message) error
	GetAllMessages(chatId string) ([]*Message, error)
}
