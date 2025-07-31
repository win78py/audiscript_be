package auth

// RegisterDTO input register
type RegisterDTO struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=3"`
}

// LoginDTO input login
type LoginDTO struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// RefreshDTO input refresh-token
type RefreshDTO struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

// APIResponse định dạng chung
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
}