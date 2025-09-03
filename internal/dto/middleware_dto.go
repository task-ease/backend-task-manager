package dto

import "backend-task-manager/internal/enums"

type RolesMiddlewareDto struct {
	Role    enums.UserRoles
	CanEdit bool
}
