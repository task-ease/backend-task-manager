package request

import (
	"github.com/google/uuid"
	"go-postgres-test/internal/types/user"
)

type CreateProject struct {
	WorkspaceId uuid.UUID `json:"workspaceId"`
	Name        string    `json:"name"`
	Prefix      string    `json:"prefix"`
}

type AddUserToProject struct {
	ProjectId uuid.UUID        `json:"projectId"`
	Role      user.ProjectRole `json:"role"`
}
