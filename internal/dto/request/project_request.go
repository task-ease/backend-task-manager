package request

import (
	"go-postgres-test/internal/enums"

	"github.com/google/uuid"
)

type CreateProject struct {
	WorkspaceId uuid.UUID `json:"workspaceId"`
	Name        string    `json:"name"`
	Prefix      string    `json:"prefix"`
}

type UserProjectManipulation struct {
	ProjectId uuid.UUID `json:"projectId"`
	UserId    uuid.UUID `json:"userId"`
}

type ChangeUserUserRoles struct {
	ProjectId uuid.UUID       `json:"projectId"`
	UserId    uuid.UUID       `json:"userId"`
	Role      enums.UserRoles `json:"role"`
}
