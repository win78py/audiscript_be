package auth

import (
	"gorm.io/gorm"
	"audiscript_be/internal/models"
)

// Repository interface
type Repository interface {
	CreateUser(user *models.User) error
	GetUserByEmail(email string) (*models.User, error)
	SaveRefreshToken(rt *models.RefreshToken) error
	GetRefreshToken(token string) (*models.RefreshToken, error)
	DeleteRefreshToken(token string) error
}

// GORM implementation
type gormRepository struct {
	db *gorm.DB
}

// NewRepository khởi tạo repo
func NewRepository(db *gorm.DB) Repository {
	db.AutoMigrate(&models.User{}, &models.RefreshToken{})
	return &gormRepository{db}
}

func (r *gormRepository) CreateUser(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *gormRepository) GetUserByEmail(email string) (*models.User, error) {
	var u models.User
	if err := r.db.Where("email = ?", email).First(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *gormRepository) SaveRefreshToken(rt *models.RefreshToken) error {
	return r.db.Create(rt).Error
}

func (r *gormRepository) GetRefreshToken(token string) (*models.RefreshToken, error) {
	var rt models.RefreshToken
	if err := r.db.Where("token = ?", token).First(&rt).Error; err != nil {
		return nil, err
	}
	return &rt, nil
}

func (r *gormRepository) DeleteRefreshToken(token string) error {
	return r.db.Where("token = ?", token).Delete(&models.RefreshToken{}).Error
}