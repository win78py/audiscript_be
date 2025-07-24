package cloudinary

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type Service interface {
	UploadAudio(ctx context.Context, file io.Reader, filename string) (string, error)
	DeleteAudio(ctx context.Context, publicID string) error
}

type service struct {
	client *CloudinaryClient
}

func NewService(client *CloudinaryClient) Service {
	return &service{client: client}
}

func (s *service) UploadAudio(ctx context.Context, file io.Reader, filename string) (string, error) {
	publicID := strings.TrimSuffix(filename, filepath.Ext(filename))
	uploadResult, err := s.client.Client.Upload.Upload(ctx, file, uploader.UploadParams{
		PublicID:     publicID,
		ResourceType: "auto",
		Folder:       "audiscript/audio",
	})
	if err != nil {
		return "", fmt.Errorf("upload failed: %w", err)
	}
	return uploadResult.SecureURL, nil
}

func (s *service) DeleteAudio(ctx context.Context, publicID string) error {
	_, err := s.client.Client.Upload.Destroy(ctx, uploader.DestroyParams{
		PublicID: publicID,
	})
	if err != nil {
		return fmt.Errorf("delete failed: %w", err)
	}
	return nil
}
