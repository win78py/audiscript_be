package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"gorm.io/gorm"
)

type StringArray []string

func (a StringArray) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func (a *StringArray) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &a)
}

type Audio struct {
	ID            string         `gorm:"type:uuid;primaryKey" json:"id"`
	Title         string         `gorm:"type:varchar(255);not null" json:"title"`
	FileURL       string         `gorm:"type:text;not null" json:"file_url"`
	Transcript    string         `gorm:"type:text" json:"transcript"`
	FileSize      int64          `gorm:"type:bigint" json:"file_size"`
	Language      string         `gorm:"type:text" json:"language"`
	Tags          StringArray    `gorm:"type:json" json:"tags"`
	CreatedBy     string         `gorm:"type:varchar(100)" json:"created_by"`
	CreatedAt     time.Time      `json:"created_at"`
	CreatedUpdate time.Time      `json:"created_update"`
	CreatedDelete gorm.DeletedAt `gorm:"index" json:"-"`
	UserID        *string        `gorm:"type:uuid;index" json:"user_id"`
	User          User           `gorm:"foreignKey:UserID" json:"-"`
}
