package domain

import (
	"context"
	"github.com/google/uuid"
	"go-postgres-test/internal/entities"
	"go-postgres-test/internal/enums"
	"io"
	"time"
)

type Message struct {
	ID             string            `db:"ID" json:"ID"`
	ChatID         string            `db:"chatID" json:"chatID"`
	SenderID       uuid.UUID         `db:"senderID" json:"senderID"`
	Content        string            `db:"content" json:"content"`
	MessageType    enums.MessageType `db:"messageType" json:"messageType"`
	CreatedAt      time.Time         `db:"createdAt" json:"createdAt"`
	UpdatedAt      time.Time         `db:"updatedAt" json:"updatedAt"`
	IsEdited       bool              `db:"isEdited" json:"isEdited"`
	IsDeleted      bool              `db:"isDeleted" json:"isDeleted"`
	IsRead         bool              `db:"isRead" json:"isRead"`
	ReplyTo        *uuid.UUID        `db:"replyTo" json:"replyTo"`
	Attachments    []Attachment      `db:"attachments" json:"attachments"`
	UnreadUsersIds []uuid.UUID       `db:"unreadUsersIds" json:"unreadUsersIds"`
}

type Attachment struct {
	ID         uuid.UUID         `json:"id"`
	MessageID  string            `json:"messageID"`
	ChatID     string            `json:"chatID"`
	FileUrl    string            `json:"fileUrl"`
	FileType   enums.MessageType `json:"fileType"`
	FileSize   int64             `json:"fileSize"`
	FileName   string            `json:"fileName"`
	Size       int64             `json:"size"`
	UploadedAt time.Time         `json:"uploadedAt"`
}

type MessageRepository interface {
	AddMessage(message *Message) (*[]uuid.UUID, error)
	UploadImage(image io.Reader) (string, error)
	GetAllMessages(chatId string) ([]*Message, error)
	AddAttachment(attachment *Attachment) error
	SetMessageRead(message *Message, userId uuid.UUID) error
	SetMessagesRead(messages *[]Message, userId uuid.UUID) error
	GetAllChatImages(chatId string) (*[]Attachment, error)

	AddMessageTx(ctx context.Context, exec entities.Execer, message *Message) (*[]uuid.UUID, error)
	GetAllMessagesTx(ctx context.Context, exec entities.Execer, chatId string) ([]*Message, error)
	GetAttachmentsTx(ctx context.Context, exec entities.Execer, messageId string, attachments *[]Attachment) error
	SetMessageReadTx(ctx context.Context, exec entities.Execer, message *Message, userId uuid.UUID) error
}
