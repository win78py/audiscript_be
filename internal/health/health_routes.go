package health

import (
	"audiscript_be/database"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Register(r gin.IRouter, db *gorm.DB) {
    svc := NewService(database.New())
    h   := NewHandler(svc)

    // GET /health
    r.GET("/", h.Check)
}