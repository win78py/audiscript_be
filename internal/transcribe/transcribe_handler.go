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

	if err := h.svc.TranscribeStream(audio, file, fileHeader.Filename); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Transcribe successful",
		"file_url": audio.FileURL,
	})
}

func (h *Handler) ListAudio(c *gin.Context) {
    audios, err := h.svc.GetAllAudio()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, audios)
}

func (h *Handler) GetAudio(c *gin.Context) {
    id := c.Param("id")
    audio, err := h.svc.GetAudioByID(id)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Audio not found"})
        return
    }
    c.JSON(http.StatusOK, audio)
}