package middleware

import (
	"backend-task-manager/internal/domain"
	"net/http"

	"github.com/gin-gonic/gin"
)

func JWTMiddleware(authService domain.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr, err := c.Cookie("token")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		userId, err := authService.VerifyToken(tokenStr)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		c.Set("userId", userId)
		c.Next()
	}

}
