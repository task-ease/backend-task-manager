package middleware

import (
	"context"
	"go-postgres-test/internal/dto"
	"go-postgres-test/internal/enums"
	"go-postgres-test/mixins"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type usecaseCheckMethod interface {
	CheckUserAccess(ctx context.Context, userId, resourceId uuid.UUID) (dto.RolesMiddlewareDto, error)
}

func AccessMiddleware(param enums.ParamKey, usecase usecaseCheckMethod, allowedRoles []enums.UserRoles) gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, err := mixins.ParseContextUserId(c)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		resourceId, err := mixins.ParamToUUID(c, string(param))

		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		data, err := usecase.CheckUserAccess(c.Request.Context(), userId, resourceId)

		if err != nil {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		if data.Role == enums.NoAccess || !mixins.Contains(allowedRoles, data.Role) {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		c.Set("canEdit", data.CanEdit)
		c.Set("role", data.Role)

		c.Next()
	}
}
