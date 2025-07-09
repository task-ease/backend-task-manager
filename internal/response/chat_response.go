package response

import (
	"go-postgres-test/internal/types"
	"go-postgres-test/internal/types/user"
	"time"
)

type GetChats struct {
	ID                    string            `json:"id"`
	Muted                 bool              `json:"muted"`
	Pinned                bool              `json:"pinned"`
	Notification          bool              `json:"notification"`
	Role                  user.UserRoles    `json:"role"`
	Type                  types.ChatType    `json:"type"`
	Name                  string            `json:"name"`
	LastMessage           string            `json:"lastMessage"`
	LastMessageTime       time.Time         `json:"lastMessageTime"`
	LastMessageType       types.MessageType `json:"lastMessageType"`
	LastMessageAttachment *string           `json:"lastMessageAttachment"`
	IsOnline              bool              `json:"isOnline"`
	ParticipantID         *string           `json:"participantId"`
	UnreadForUser         int               `json:"unreadForUser"`
}

type GetChatsSearch struct {
	ID          *string `json:"id"`
	UserID      *string `json:"userId"`
	Name        string  `json:"name"`
	LastMessage *string `json:"lastMessage"`
}
