package request

import (
	"go-postgres-test/internal/enums"

	"github.com/google/uuid"
)

type CreateDocsRequest struct {
	Name       string                        `json:"name"`
	ParentId   *uuid.UUID                    `json:"parentId"`
	ProjectId  *uuid.UUID                    `json:"projectId"`
	Visibility enums.DocumentVisibilityTypes `json:"visibility"`
}

type UpdateDocUserEditPermissionsRequest struct {
	To     bool      `json:"to"`
	UserId uuid.UUID `json:"userId"`
}
