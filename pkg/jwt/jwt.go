package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
	"audiscript_be/config"
	"errors"
)

// Claims chứa thông tin mã hóa trong JWT
type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

// GenerateTokens trả về access & refresh token
func GenerateTokens(userID string) (accessToken, refreshToken string, err error) {
	cfg := config.AppConfig.JWT

	// Access Token
	atClaims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(cfg.AccessExpiry) * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	accessToken, err = at.SignedString([]byte(cfg.Secret))
	if err != nil {
		return
	}

	// Refresh Token
	rtClaims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(cfg.RefreshExpiry) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	refreshToken, err = rt.SignedString([]byte(cfg.Secret))
	return
}

// VerifyToken kiểm tra và trích thông tin từ token
func VerifyToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.AppConfig.JWT.Secret), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}