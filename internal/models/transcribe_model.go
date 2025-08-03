package models

import (
	"time"

	"gorm.io/gorm"
)

type Audio struct {
	ID            string         `gorm:"type:uuid;primaryKey" json:"id"`
	Title         string         `gorm:"type:varchar(255);not null" json:"title"`
	FileURL       string         `gorm:"type:text;not null" json:"file_url"`
	Transcript    string         `gorm:"type:text" json:"transcript"`
	CreatedBy     string         `gorm:"type:varchar(100)" json:"created_by"`
	CreatedAt     time.Time      `json:"created_at"`
	CreatedUpdate time.Time      `json:"created_update"`
	CreatedDelete gorm.DeletedAt `gorm:"index" json:"-"`
	UserID        uint           `gorm:"not null;index" json:"user_id"`
	User          User      `gorm:"foreignKey:UserID" json:"-"` // Quan hệ với User|

}
