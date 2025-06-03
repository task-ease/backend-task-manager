package middleware

import (
	"github.com/gin-gonic/gin"
	"go-postgres-test/internal/domain"
	"net/http"
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
