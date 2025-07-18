package health

import (
    "net/http"

    "github.com/gin-gonic/gin"
)

type Handler struct {
    svc Service
}

func NewHandler(svc Service) *Handler {
    return &Handler{svc: svc}
}

func (h *Handler) Check(c *gin.Context) {
    resp := h.svc.Check()
    c.JSON(http.StatusOK, resp)
}