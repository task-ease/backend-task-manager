package response

import (
	"backend-task-manager/internal/enums"
	"time"
)

type GetChats struct {
	ID                    string            `json:"id"`
	Muted                 bool              `json:"muted"`
	Pinned                bool              `json:"pinned"`
	Notification          bool              `json:"notification"`
	Role                  enums.UserRoles   `json:"role"`
	Type                  enums.ChatType    `json:"type"`
	Name                  string            `json:"name"`
	LastMessage           string            `json:"lastMessage"`
	LastMessageTime       time.Time         `json:"lastMessageTime"`
	LastMessageType       enums.MessageType `json:"lastMessageType"`
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
