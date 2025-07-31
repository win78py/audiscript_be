package auth

import "gorm.io/gorm"

// Repository interface
type Repository interface {
	CreateUser(user *User) error
	GetUserByEmail(email string) (*User, error)
	SaveRefreshToken(rt *RefreshToken) error
	GetRefreshToken(token string) (*RefreshToken, error)
	DeleteRefreshToken(token string) error
}

// GORM implementation
type gormRepository struct {
	db *gorm.DB
}

// NewRepository khởi tạo repo
func NewRepository(db *gorm.DB) Repository {
	db.AutoMigrate(&User{}, &RefreshToken{})
	return &gormRepository{db}
}

func (r *gormRepository) CreateUser(user *User) error {
	return r.db.Create(user).Error
}

func (r *gormRepository) GetUserByEmail(email string) (*User, error) {
	var u User
	if err := r.db.Where("email = ?", email).First(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *gormRepository) SaveRefreshToken(rt *RefreshToken) error {
	return r.db.Create(rt).Error
}

func (r *gormRepository) GetRefreshToken(token string) (*RefreshToken, error) {
	var rt RefreshToken
	if err := r.db.Where("token = ?", token).First(&rt).Error; err != nil {
		return nil, err
	}
	return &rt, nil
}

func (r *gormRepository) DeleteRefreshToken(token string) error {
	return r.db.Where("token = ?", token).Delete(&RefreshToken{}).Error
}