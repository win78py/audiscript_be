package transcribe

import (
	// "audiscript_be/config"
	// "audiscript_be/pkg/jwt"

	"github.com/gin-gonic/gin"
)

func Register(r *gin.Engine, svc Service) {
    h := NewHandler(svc)
    audioGroup := r.Group("/audio")
    audioGroup.POST("/create", h.CreateAudio)
    audioGroup.POST("/transcribe", h.Transcribe)
    audioGroup.GET("/", h.ListAudio)
    // audioGroup.GET("/", jwt.AuthGuard(config.AppConfig.JWT.Secret), h.ListAudio)
    audioGroup.GET("/:id", h.GetAudio)
}