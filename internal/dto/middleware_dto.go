package dto

import "go-postgres-test/internal/enums"

type RolesMiddlewareDto struct {
	Role    enums.UserRoles
	CanEdit bool
}
