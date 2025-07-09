package domain

import (
	"github.com/google/uuid"
	"go-postgres-test/internal/types"
	"io"
	"time"
)

type Message struct {
	ID             string            `db:"ID" json:"ID"`
	ChatID         string            `db:"chatID" json:"chatID"`
	SenderID       uuid.UUID         `db:"senderID" json:"senderID"`
	Content        string            `db:"content" json:"content"`
	MessageType    types.MessageType `db:"messageType" json:"messageType"`
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
	FileType   types.MessageType `json:"fileType"`
	FileSize   int64             `json:"fileSize"`
	FileName   string            `json:"fileName"`
	Size       int64             `json:"size"`
	UploadedAt time.Time         `json:"uploadedAt"`
}

type MessageRepository interface {
	AddMessage(message *Message) error
	UploadImage(image io.Reader) (string, error)
	GetAllMessages(chatId string) ([]*Message, error)
	AddAttachment(attachment *Attachment) error
	SetMessageRead(message *Message, userId uuid.UUID) (*string, error)
	SetMessagesRead(messages *[]Message, userId uuid.UUID) (*[]string, error)
}
