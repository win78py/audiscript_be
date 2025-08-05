package transcribe

import (
	// "audiscript_be/internal/auth"
	"audiscript_be/internal/models"
	"audiscript_be/pkg/pagination"
	"context"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	svc Service
	// authService auth.Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) CreateAudio(c *gin.Context) {
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
	userID := c.PostForm("user_id")
	var userIDPtr *string
	if userID != "" && userID != "undefined" {
		userIDPtr = &userID
	} else {
		userIDPtr = nil
	}
	tags := c.PostFormArray("tags")
	defer file.Close()

	audio := &models.Audio{
		ID:            uuid.New().String(),
		Title:         fileHeader.Filename,
		FileURL:       "",
		Transcript:    "",
		FileSize:      fileHeader.Size,
		Language:      models.AutoLanguageDetection,
		Tags:          tags,
		CreatedAt:     time.Now(),
		CreatedUpdate: time.Now(),
		UserID:        userIDPtr,
	}

	log.Printf("Uploading audio (stream): %s", fileHeader.Filename)

	if err := h.svc.CreateAudio(audio, file, fileHeader.Filename, fileHeader.Size); err != nil {
		log.Printf("Handler error: %v", err)
		if errors.Is(err, context.DeadlineExceeded) || strings.Contains(err.Error(), "timeout") {
			c.JSON(http.StatusGatewayTimeout, gin.H{
				"error":   "Create timeout",
				"details": err.Error(),
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	// user, err := h.authService.GetByID(*audio.UserID)
	// if err != nil {
	// 	c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
	// 	return
	// }
	// email := user.Email

	c.JSON(http.StatusOK, gin.H{
		"message":   "Transcribe successful",
		"id":        audio.ID,
		"title":     audio.Title,
		"file_url":  audio.FileURL,
		"file_size": audio.FileSize,
		"language":  audio.Language,
		"tags":      audio.Tags,
		"user_id":   audio.UserID,
		// "email":    email,
	})
}

func (h *Handler) Transcribe(c *gin.Context) {
    var req TranscribeRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
        return
    }

    transcript, err := h.svc.Transcribe(req.FileURL, req.Language)
    if err != nil {
        if errors.Is(err, context.DeadlineExceeded) || strings.Contains(err.Error(), "timeout") {
            c.JSON(http.StatusGatewayTimeout, gin.H{
                "error":   "Transcribe timeout",
                "details": err.Error(),
            })
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        }
        return
    }

    if err := h.svc.UpdateTranscript(req.AudioID, transcript, req.Language); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update transcript"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "transcript": transcript,
        "language":   req.Language,
    })
}

func (h *Handler) ListAudio(c *gin.Context) {
	var pageReq pagination.PageRequest
	if err := c.ShouldBindQuery(&pageReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pagination parameters"})
		return
	}

	userID := c.Query("user_id")
	var userIDPtr *string
	if userID != "" {
		userIDPtr = &userID
	}
	log.Printf("Listing audio for user: %s, page: %d, limit: %d", userID, pageReq.Page, pageReq.Limit)

	result, err := h.svc.ListAudio(c.Request.Context(), pageReq.Page, pageReq.Limit, userIDPtr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
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
