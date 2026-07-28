package request

import "github.com/google/uuid"

type AddUserToChat struct {
	ChatID      string    `json:"chatId"`
	UserID      uuid.UUID `json:"userId"`
	WorkspaceID uuid.UUID `json:"workspaceId"`
}
