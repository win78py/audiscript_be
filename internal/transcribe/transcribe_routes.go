package transcribe

import (
	"audiscript_be/internal/app"

	"github.com/gin-gonic/gin"
)

func Register(r *gin.RouterGroup, deps *app.AppDependencies) {
    repo := NewRepository(deps.DB)
    svc := NewService(repo, deps.Cloudinary)
    h := NewHandler(svc)

    r.POST("/transcribe", h.Transcribe)
    r.GET("/", h.ListAudio)
    r.GET("/:id", h.GetAudio)
}