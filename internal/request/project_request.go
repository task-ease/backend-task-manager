package request

import "github.com/google/uuid"

type CreateProject struct {
	WorkspaceId uuid.UUID `json:"workspaceId"`
	Name        string    `json:"name"`
	Prefix      string    `json:"prefix"`
}
