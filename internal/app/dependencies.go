package app

import (
	"audiscript_be/config"
	"audiscript_be/database"
	"audiscript_be/internal/auth"
	"audiscript_be/internal/cloudinary"
	"audiscript_be/internal/transcribe"
	"log"

	"gorm.io/gorm"
)

type AppDependencies struct {
    DB         *gorm.DB
    Cloudinary cloudinary.Service
    AuthService     auth.Service
    TranscribeService transcribe.Service
    // Redis     redis.Client 
}

func NewDependencies() *AppDependencies {
    // config.LoadConfig()

    dbSvc := database.New()
    db := dbSvc.DB()

    cldClient, err := cloudinary.NewClient(config.AppConfig.Cloudinary)
    if err != nil {
        log.Fatalf("Failed to create Cloudinary client: %v", err)
    }
    cldSvc := cloudinary.NewService(cldClient)

    // Khởi repo + svc cho auth
    authRepo := auth.NewRepository(db)
    authSvc  := auth.NewService(authRepo)

    // Khởi repo + svc cho transcribe
    transRepo := transcribe.NewRepository(db)
    transSvc  := transcribe.NewService(transRepo, cldSvc)

    return &AppDependencies{
        DB:         db,
        Cloudinary: cldSvc,
        AuthService:       authSvc,
        TranscribeService: transSvc,
    }
}