package transcribe

import (
	"context"

	"gorm.io/gorm"
    "audiscript_be/internal/models"
)

type Repository interface {
	Save(ctx context.Context, t *models.Audio) error
	GetAll(ctx context.Context) ([]models.Audio, error)
    GetByID(ctx context.Context, id string) (*models.Audio, error)
    GetPaginated(ctx context.Context, filters map[string]interface{}, offset, limit int) ([]models.Audio, error)
	Count(ctx context.Context, filters map[string]interface{}) (int64, error)
}

type repo struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	// Auto migrate
	db.AutoMigrate(&models.Audio{})
	return &repo{db: db}
}

func (r *repo) Save(ctx context.Context, t *models.Audio) error {
	return r.db.WithContext(ctx).Create(t).Error
}

func (r *repo) GetAll(ctx context.Context) ([]models.Audio, error) {
    var audios []models.Audio
    err := r.db.WithContext(ctx).Find(&audios).Error
    return audios, err
}

func (r *repo) GetByID(ctx context.Context, id string) (*models.Audio, error) {
    var audio models.Audio
    err := r.db.WithContext(ctx).First(&audio, "id = ?", id).Error
    if err != nil {
        return nil, err
    }
    return &audio, nil
}

func (r *repo) GetPaginated(ctx context.Context, filters map[string]interface{}, offset, limit int) ([]models.Audio, error) {
	var audios []models.Audio
	query := r.db.WithContext(ctx).Model(&models.Audio{})

	for k, v := range filters {
		query = query.Where(k+" = ?", v)
	}

	err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&audios).Error
	return audios, err
}

func (r *repo) Count(ctx context.Context, filters map[string]interface{}) (int64, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&models.Audio{})

	for k, v := range filters {
		query = query.Where(k+" = ?", v)
	}

	err := query.Count(&count).Error
	return count, err
}