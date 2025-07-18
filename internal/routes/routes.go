package routes

import (

	"github.com/gin-gonic/gin"

	"audiscript_be/internal/app"
	"audiscript_be/internal/health"
	"audiscript_be/internal/transcribe"
)

func RegisterAll(r *gin.Engine, deps *app.AppDependencies) {
    // CORS, Recovery, Logging… nếu cần
    audioGroup := r.Group("/audio")
    transcribe.Register(audioGroup, deps)

	healthGroup := r.Group("/health")
    health.Register(healthGroup, deps.DB)
}
