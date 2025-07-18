package app

import (
	"audiscript_be/internal/cloudinary"
	"gorm.io/gorm"
)

type AppDependencies struct {
	DB        *gorm.DB
	Cloudinary cloudinary.Service
}
