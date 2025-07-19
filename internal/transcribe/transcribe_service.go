package transcribe

import (
	"audiscript_be/internal/cloudinary"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type Service interface {
	TranscribeStream(t *Audio, file io.Reader, filename string) error
	GetAllAudio() ([]Audio, error)
	GetAudioByID(id string) (*Audio, error)
}


type service struct {
    repo Repository
    cld  cloudinary.Service
}

func NewService(r Repository, c cloudinary.Service) Service {
    return &service{repo: r, cld: c}
}

func (s *service) TranscribeStream(t *Audio, file io.Reader, filename string) error {
    log.Printf("Transcribing audio: %s", filename)

    // 1. Upload file lên Cloudinary
    url, err := s.cld.UploadAudio(context.Background(), file, filename)
    if err != nil {
        return err
    }
    t.FileURL = url
    log.Printf("Uploaded audio to Cloudinary: %s", url)

    // 2. Gọi HTTP tới Python service để transcribe
    transcript, err := s.callPythonTranscribe(url)
    if err != nil {
        log.Printf("Transcription failed: %v", err)
        t.Transcript = ""
    } else {
        t.Transcript = transcript
    }

    // 3. Lưu vào database
    return s.repo.Save(context.Background(), t)
}

func (s *service) callPythonTranscribe(audioURL string) (string, error) {
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

    log.Printf("Calling Python service at: %s", pyServiceURL)

    resp, err := http.Post(pyServiceURL, "application/json", bytes.NewBuffer(reqBody))
    if err != nil {
        return "", fmt.Errorf("failed to call python service: %w", err)
    }
    defer resp.Body.Close()

    var result struct {
        Transcript string `json:"transcript"`
        Error      string `json:"error"`
    }
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return "", fmt.Errorf("failed to decode python service response: %w", err)
    }
    if result.Error != "" {
        return "", fmt.Errorf("python service error: %s", result.Error)
    }
    return result.Transcript, nil
}

func (s *service) GetAllAudio() ([]Audio, error) {
    return s.repo.GetAll(context.Background())
}

func (s *service) GetAudioByID(id string) (*Audio, error) {
    return s.repo.GetByID(context.Background(), id)
}