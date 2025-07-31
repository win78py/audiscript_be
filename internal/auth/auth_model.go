package auth

import (
    "time"
    "gorm.io/gorm"
)

// User represents application user
type User struct {
    ID        uint           `gorm:"primaryKey" json:"id"`
    Email     string         `gorm:"unique;not null" json:"email"`
    Password  string         `gorm:"not null" json:"-"`
    CreatedAt time.Time      `json:"createdAt"`
    UpdatedAt time.Time      `json:"updatedAt"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// RefreshToken to support token revocation
type RefreshToken struct {
    ID        uint      `gorm:"primaryKey"`
    Token     string    `gorm:"unique;not null"`
    UserID    uint      `gorm:"index;not null"`
    ExpiresAt time.Time `gorm:"not null"`
    CreatedAt time.Time
}