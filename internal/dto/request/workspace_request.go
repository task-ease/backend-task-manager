package request

import (
	"go-postgres-test/internal/enums"

	"github.com/google/uuid"
)

type AddUserToWorkspace struct {
	WorkSpaceId uuid.UUID       `json:"workSpaceId"`
	UserId      uuid.UUID       `json:"userId"`
	Role        enums.UserRoles `json:"role"`
}

type WorkspaceUserManipulation struct {
	UserId      uuid.UUID `json:"userId"`
	WorkspaceId uuid.UUID `json:"workSpaceId"`
}

type WorkspaceUserManipulationRole struct {
	UserId      uuid.UUID       `json:"userId"`
	WorkspaceId uuid.UUID       `json:"workSpaceId"`
	Role        enums.UserRoles `json:"role"`
}
