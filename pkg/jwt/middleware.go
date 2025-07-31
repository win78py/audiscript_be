package jwt

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthGuard bảo vệ route yêu cầu JWT hợp lệ
func AuthGuard(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		parts := strings.Fields(header)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "invalid or missing auth header"})
			return
		}

		claims, err := VerifyToken(parts[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
			return
		}
		c.Set("currentUser", claims.UserID)
		c.Next()
	}
}