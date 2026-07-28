package entities

import (
	"backend-task-manager/internal/enums"
	"time"

	"github.com/google/uuid"
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
