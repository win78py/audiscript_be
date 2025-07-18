package cloudinary

import "os"

type Config struct {
	CloudName string
	APIKey    string
	APISecret string
}

func LoadConfig() Config {
	return Config{
		CloudName: os.Getenv("CLOUDINARY_CLOUD_NAME"),
		APIKey:    os.Getenv("CLOUDINARY_API_KEY"),
		APISecret: os.Getenv("CLOUDINARY_API_SECRET"),
	}
}
