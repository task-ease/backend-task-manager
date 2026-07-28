package request

import (
	"backend-task-manager/internal/enums"

	"github.com/google/uuid"
)

type CreateProject struct {
	Name   string `json:"name"`
	Prefix string `json:"prefix"`
}

type UserId struct {
	UserId uuid.UUID `json:"userId"`
}

type ChangeUserUserRoles struct {
	UserId uuid.UUID       `json:"userId"`
	Role   enums.UserRoles `json:"role"`
}
