package transcribe

import (
	"context"

	"gorm.io/gorm"
)

type Repository interface {
	Save(ctx context.Context, t *Audio) error
	GetAll(ctx context.Context) ([]Audio, error)
    GetByID(ctx context.Context, id string) (*Audio, error)
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

func (r *repo) GetAll(ctx context.Context) ([]Audio, error) {
    var audios []Audio
    err := r.db.WithContext(ctx).Find(&audios).Error
    return audios, err
}

func (r *repo) GetByID(ctx context.Context, id string) (*Audio, error) {
    var audio Audio
    err := r.db.WithContext(ctx).First(&audio, "id = ?", id).Error
    if err != nil {
        return nil, err
    }
    return &audio, nil
}