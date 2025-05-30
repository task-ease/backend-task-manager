package middleware

import (
	"github.com/gin-gonic/gin"
	"go-postgres-test/internal/domain"
)

func AuthMiddleware(authService domain.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}
