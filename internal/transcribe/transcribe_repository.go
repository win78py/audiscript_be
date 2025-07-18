package transcribe

import (
	"context"

	"gorm.io/gorm"
)

type Repository interface {
	Save(ctx context.Context, t *Audio) error
}

type repo struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	// Auto migrate
	db.AutoMigrate(&Audio{})
	return &repo{db: db}
}

func (r *repo) Save(ctx context.Context, t *Audio) error {
	return r.db.WithContext(ctx).Create(t).Error
}
