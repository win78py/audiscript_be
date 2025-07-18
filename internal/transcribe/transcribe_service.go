package transcribe

import (
	"audiscript_be/internal/cloudinary"
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os/exec"
	"strings"
)

type Service interface {
	TranscribeStream(t *Audio, file io.Reader, filename string) error
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
    
    // 2. Gọi Python script để transcribe
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
    cmd := exec.Command("python", "internal/transcribe/transcribe.py", audioURL)
    var out bytes.Buffer
    cmd.Stdout = &out
    cmd.Stderr = nil
    err := cmd.Run()
    if err != nil {
        return "", fmt.Errorf("failed to run python script: %w", err)
    }
    return strings.TrimSpace(out.String()), nil
    
}