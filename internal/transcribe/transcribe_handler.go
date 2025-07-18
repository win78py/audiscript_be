package transcribe

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	svc Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Transcribe(c *gin.Context) {
	title := c.PostForm("title")

	fileHeader, err := c.FormFile("file_url")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
		return
	}

	// Mở file trực tiếp từ multipart (stream)
	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not open uploaded file"})
		return
	}
	defer file.Close()

	audio := &Audio{
		ID:            uuid.New().String(),
		Title:         title,
		FileURL:       "",
		CreatedAt:     time.Now(),
		CreatedUpdate: time.Now(),
	}

	log.Printf("Uploading audio (stream): %s", fileHeader.Filename)

	// Gọi service: truyền stream và filename
	if err := h.svc.TranscribeStream(audio, file, fileHeader.Filename); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Transcribe successful",
		"file_url": audio.FileURL,
	})
}
