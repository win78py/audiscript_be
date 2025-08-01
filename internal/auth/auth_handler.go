package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Handler kết nối HTTP <-> Service
type Handler struct {
	service Service
}

func NewHandler(s Service) *Handler {
	return &Handler{service: s}
}

// Register endpoint
func (h *Handler) Register(c *gin.Context) {
	var dto RegisterDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{Success: false, Message: err.Error()})
		return
	}
	user, err := h.service.Register(dto.Email, dto.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{Success: false, Message: err.Error()})
		return
	}
	c.JSON(http.StatusCreated, APIResponse{Success: true, Data: user, Message: "[register] successfully"})
}

// Login endpoint
func (h *Handler) Login(c *gin.Context) {
	var dto LoginDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{Success: false, Message: err.Error()})
		return
	}
	access, refresh, user, err := h.service.Login(dto.Email, dto.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, APIResponse{Success: false, Message: err.Error()})
		return
	}
	// thêm ip & permission nếu cần
	resp := map[string]interface{}{  
		"token":        access,
		"refreshToken": refresh,
		"customer": map[string]interface{}{  
			"id":    user.ID,
			"email": user.Email,
			"ip":    c.ClientIP(),
		},
	}
	c.JSON(http.StatusOK, APIResponse{Success: true, Data: resp, Message: "[login] login successfully."})
}

// Refresh endpoint
func (h *Handler) Refresh(c *gin.Context) {
	var dto RefreshDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{Success: false, Message: err.Error()})
		return
	}
	newAt, newRt, err := h.service.Refresh(dto.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, APIResponse{Success: false, Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, APIResponse{Success: true, Data: map[string]string{"token": newAt, "refreshToken": newRt}, Message: "[refresh] refresh token successfully."})
}