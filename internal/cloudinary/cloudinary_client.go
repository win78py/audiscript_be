package cloudinary

import (
	"fmt"

	cld "github.com/cloudinary/cloudinary-go/v2"
	// "github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type CloudinaryClient struct {
	Client *cld.Cloudinary
}

func NewClient(cfg Config) (*CloudinaryClient, error) {
	cldClient, err := cld.NewFromParams(cfg.CloudName, cfg.APIKey, cfg.APISecret)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize cloudinary client: %w", err)
	}
	return &CloudinaryClient{Client: cldClient}, nil
}
