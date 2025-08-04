package transcribe

import (
	"audiscript_be/internal/cloudinary"
	"audiscript_be/internal/models"
	"audiscript_be/pkg/pagination"
	"audiscript_be/pkg/util"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type Service interface {
	TranscribeStream(t *models.Audio, file io.Reader, filename string, fileSize int64) error
	GetAudioByID(id string) (*models.Audio, error)
	ListAudio(ctx context.Context, page, limit int, userID *string) (*pagination.PageResponse[models.Audio], error)
}

type service struct {
	repo Repository
	cld  cloudinary.Service
}

func NewService(r Repository, c cloudinary.Service) Service {
	return &service{repo: r, cld: c}
}

func (s *service) TranscribeStream(t *models.Audio, file io.Reader, filename string, fileSize int64) error {
	log.Printf("Transcribing audio: %s", filename)

	// 1. Upload file lên Cloudinary
	url, err := s.cld.UploadAudio(context.Background(), file, filename)
	if err != nil {
		return err
	}
	t.FileURL = url
	log.Printf("Uploaded audio to Cloudinary: %s", url)

	// 2. Gọi HTTP tới Python service để transcribe
	transcript, err := s.callPythonTranscribe(url, fileSize)
	if err != nil {
		log.Printf("Transcription failed: %v", err)
		t.Transcript = ""
		return err
	} else {
		t.Transcript = transcript
	}

	// 3. Lưu vào database
	return s.repo.Save(context.Background(), t)
}

func (s *service) callPythonTranscribe(audioURL string, fileSize int64) (string, error) {
	// Chuẩn bị body JSON
	reqBody, _ := json.Marshal(map[string]string{
		"file_url": audioURL,
	})

	pyServiceURL := os.Getenv("PY_SERVICE_URL")
	if pyServiceURL == "" {
		pyServiceURL = "http://localhost:8000"
	}

	if pyServiceURL[len(pyServiceURL)-len("/transcribe"):] != "/transcribe" {
		pyServiceURL = pyServiceURL + "/transcribe"
	}

	// log.Printf("Calling Python service at: %s", pyServiceURL)
	// Tạo HTTP Client với timeout
	client := util.DefaultHTTPClient

	// var timeout time.Duration
	// switch {
	// case fileSize > 8000*1024:
	//     timeout = 999 * time.Second
	// case fileSize > 5000*1024:
	//     timeout = 80 * time.Second
	// case fileSize > 2000*1024:
	//     timeout = 60 * time.Second
	// default:
	//     timeout = 30 * time.Second
	// }
	// Gửi request có context timeout

	var timeout time.Duration = 5 * time.Minute
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	start := time.Now()

	req, err := http.NewRequestWithContext(ctx, "POST", pyServiceURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	elapsed := time.Since(start)
	log.Printf("callPythonTranscribe: request took %v", elapsed)
	if err != nil {
		return "", fmt.Errorf("failed to call python service: %w", err)
	}
	defer resp.Body.Close()

	var result struct {
		Transcript string `json:"transcript"`
		Error      string `json:"error"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	if result.Error != "" {
		return "", fmt.Errorf("python service error: %s", result.Error)
	}
	return result.Transcript, nil
}

func (s *service) ListAudio(ctx context.Context, page, limit int, userID *string) (*pagination.PageResponse[models.Audio], error) {
	offset := pagination.GetOffset(page, limit)

	filters := map[string]interface{}{}
	if userID != nil {
		filters["user_id"] = *userID
	}

	audios, err := s.repo.GetPaginated(ctx, filters, offset, limit)
	if err != nil {
		return nil, err
	}

	total, err := s.repo.Count(ctx, filters)
	if err != nil {
		return nil, err
	}

	meta := pagination.GetMetadata(total, page, limit)
	return &pagination.PageResponse[models.Audio]{
		Data:     audios,
		Metadata: meta,
	}, nil
}

func (s *service) GetAudioByID(id string) (*models.Audio, error) {
	return s.repo.GetByID(context.Background(), id)
}
