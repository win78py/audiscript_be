package models

import (
	"time"

	"gorm.io/gorm"
)

// User represents application user
type User struct {
	ID           string         `gorm:"type:uuid;primaryKey" json:"id"`
	Username     string         `json:"username"`
	FirstName    string         `json:"firstName"`
	LastName     string         `json:"lastName"`
	ProfileImage string         `gorm:"default:'https://res.cloudinary.com/dj2krdujc/image/upload/v1754225216/audiscript/profile/Rose-Blackpink_trxiwz.jpg'" json:"profileImage"`
	Email        string         `gorm:"unique;not null" json:"email"`
	Password     string         `gorm:"not null" json:"-"`
	CreatedAt    time.Time      `json:"createdAt"`
	UpdatedAt    time.Time      `json:"updatedAt"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
	Audios       []Audio        `gorm:"foreignKey:UserID" json:"audios,omitempty"`
}

// RefreshToken to support token revocation
type RefreshToken struct {
	ID        uint      `gorm:"primaryKey"`
	Token     string    `gorm:"unique;not null"`
	UserID    string      `gorm:"index;not null"`
	ExpiresAt time.Time `gorm:"not null"`
	CreatedAt time.Time
}
