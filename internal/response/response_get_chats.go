package response

import (
	"go-postgres-test/internal/types"
	"time"
)

type GetChats struct {
	ID              string             `json:"id"`
	Muted           bool               `json:"muted"`
	Pinned          bool               `json:"pinned"`
	Notification    bool               `json:"notification"`
	Role            types.UserChatRole `json:"role"`
	Type            types.ChatType     `json:"type"`
	Name            string             `json:"name"`
	LastMessage     string             `json:"lastMessage"`
	LastMessageTime time.Time          `json:"lastMessageTime"`
	LastMessageType types.MessageType  `json:"lastMessageType"`
	IsOnline        bool               `json:"isOnline"`
	ParticipantID   *string            `json:"participantId"`
}
