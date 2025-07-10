package domain

import (
	"github.com/google/uuid"
	"go-postgres-test/internal/response"
	"go-postgres-test/internal/types"
	"go-postgres-test/internal/types/user"
	"time"
)

type Chat struct {
	ID              string         `json:"id"`
	WorkspaceID     uuid.UUID      `json:"workspaceId"`
	CreatorID       uuid.UUID      `json:"creatorId"`
	Type            types.ChatType `json:"type"`
	CreatedAt       time.Time      `json:"createdAt"`
	UpdatedAt       time.Time      `json:"updatedAt"`
	LastMessageTime time.Time      `json:"lastMessageTime"`
}

type GroupChat struct {
	Chat
	Name    string  `json:"name"`
	IconUrl *string `json:"iconUrl"`
}

type UserChat struct {
	ChatID       uuid.UUID     `db:"chat_id"`
	UserID       uuid.UUID     `db:"user_id"`
	Muted        bool          `db:"muted"`
	Pinned       bool          `db:"pinned"`
	Notification bool          `db:"notification"`
	Role         user.ChatRole `db:"role"`
	JoinedAt     time.Time     `db:"joined_at"`
}

type ChatRepository interface {
	CreateChat(chat *Chat, participantId uuid.UUID) error
	AddUserToChat(userId uuid.UUID, chatId string, workspaceId uuid.UUID, role user.ChatRole) error
	GetAllUserChats(userId uuid.UUID, workspaceId uuid.UUID) ([]response.GetChats, error)
	GetChatsBySearch(userID uuid.UUID, workspaceId uuid.UUID, value string) ([]response.GetChatsSearch, error)
}
