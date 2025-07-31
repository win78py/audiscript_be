package routes

import (
	"github.com/gin-gonic/gin"

	"audiscript_be/internal/app"
	"audiscript_be/internal/auth"
	"audiscript_be/internal/health"
	"audiscript_be/internal/transcribe"
)

func RegisterAll(r *gin.Engine, deps *app.AppDependencies) {
    // CORS, Recovery, Logging… nếu cần
	auth.Register(r, deps.AuthService)
    transcribe.Register(r, deps.TranscribeService)
    health.Register(r, deps.DB)
}
