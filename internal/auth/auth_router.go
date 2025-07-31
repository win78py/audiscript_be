package auth

import (
	"github.com/gin-gonic/gin"
)

// Register đăng ký các đường dẫn auth
func Register(r *gin.Engine, svc Service) {
	h := NewHandler(svc)
	// Group /auth
	authGroup := r.Group("/auth")
	authGroup.POST("/register", h.Register)
	authGroup.POST("/login", h.Login)
	authGroup.POST("/refresh", h.Refresh)
}