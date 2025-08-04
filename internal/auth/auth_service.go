package auth

import (
	"errors"
	"time"

	"audiscript_be/config"
	"audiscript_be/internal/models"
	"audiscript_be/pkg/hash"
	"audiscript_be/pkg/jwt"

	"github.com/google/uuid"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrInvalidRefreshToken = errors.New("invalid refresh token")
)

// Service interface
type Service interface {
	Register(email, password string) (*models.User, error)
	Login(email, password string) (accessToken, refreshToken string, user *models.User, err error)
	Refresh(oldToken string) (newAccess, newRefresh string, err error)
	GetByID(id string) (*models.User, error)
}

// impl
type authService struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &authService{repo}
}

func (s *authService) Register(email, password string) (*models.User, error) {
	hashPwd, err := hash.HashPassword(password)
	if err != nil {
		return nil, err
	}
	user := &models.User{
        ID:       uuid.New().String(),
        Email:    email,
        Password: hashPwd,
    }
	if err := s.repo.CreateUser(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *authService) Login(email, password string) (string, string, *models.User, error) {
	user, err := s.repo.GetUserByEmail(email)
	if err != nil {
		return "", "", nil, ErrInvalidCredentials
	}
	if !hash.CheckPasswordHash(password, user.Password) {
		return "", "", nil, ErrInvalidCredentials
	}
	// gen tokens
	access, refresh, err := jwt.GenerateTokens(user.ID)
	if err != nil {
		return "", "", nil, err
	}
	// save refresh to DB
	exp := time.Now().Add(time.Duration(config.AppConfig.JWT.RefreshExpiry) * time.Hour)
	rt := &models.RefreshToken{Token: refresh, UserID: user.ID, ExpiresAt: exp}
	_ = s.repo.SaveRefreshToken(rt)
	return access, refresh, user, nil
}

func (s *authService) Refresh(oldToken string) (string, string, error) {
	claims, err := jwt.VerifyToken(oldToken)
	if err != nil {
		return "", "", ErrInvalidRefreshToken
	}
	rt, err := s.repo.GetRefreshToken(oldToken)
	if err != nil || rt.ExpiresAt.Before(time.Now()) {
		return "", "", ErrInvalidRefreshToken
	}
	// tạo lại
	newAccess, newRefresh, err := jwt.GenerateTokens(claims.UserID)
	if err != nil {
		return "", "", err
	}
	// revoke old, save new
	_ = s.repo.DeleteRefreshToken(oldToken)
	exp := time.Now().Add(time.Duration(config.AppConfig.JWT.RefreshExpiry) * time.Hour)
	_ = s.repo.SaveRefreshToken(&models.RefreshToken{Token: newRefresh, UserID: claims.UserID, ExpiresAt: exp})
	return newAccess, newRefresh, nil
}

func (s *authService) GetByID(id string) (*models.User, error) {
    return s.repo.GetUserByID(id)
}