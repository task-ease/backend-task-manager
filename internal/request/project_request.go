package request

import (
	"github.com/google/uuid"
	"go-postgres-test/internal/enums"
)

type CreateProject struct {
	WorkspaceId uuid.UUID `json:"workspaceId"`
	Name        string    `json:"name"`
	Prefix      string    `json:"prefix"`
}

type AddUserToProject struct {
	ProjectId uuid.UUID `json:"projectId"`
	UserId    uuid.UUID `json:"userId"`
}

type ChangeUserProjectRole struct {
	ProjectId uuid.UUID         `json:"projectId"`
	UserId    uuid.UUID         `json:"userId"`
	Role      enums.ProjectRole `json:"role"`
}
