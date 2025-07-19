package config

import (
    "log"
    "os"

    "github.com/joho/godotenv"
)

type DBConfig struct {
    Host     string
    Port     string
    User     string
    Password string
    Name     string
    Schema   string
    SSLMode  string
}

type CloudinaryConfig struct {
    CloudName string
    APIKey    string
    APISecret string
}

type Config struct {
    DB         DBConfig
    Cloudinary CloudinaryConfig
}

var AppConfig *Config

func LoadConfig() {
    _ = godotenv.Load()

    AppConfig = &Config{
        DB: DBConfig{
            Host:     os.Getenv("BLUEPRINT_DB_HOST"),
            Port:     os.Getenv("BLUEPRINT_DB_PORT"),
            User:     os.Getenv("BLUEPRINT_DB_USERNAME"),
            Password: os.Getenv("BLUEPRINT_DB_PASSWORD"),
            Name:     os.Getenv("BLUEPRINT_DB_DATABASE"),
            Schema:   os.Getenv("BLUEPRINT_DB_SCHEMA"),
            SSLMode:  os.Getenv("BLUEPRINT_DB_SSLMODE"),
        },
        Cloudinary: CloudinaryConfig{
            CloudName: os.Getenv("CLOUDINARY_CLOUD_NAME"),
            APIKey:    os.Getenv("CLOUDINARY_API_KEY"),
            APISecret: os.Getenv("CLOUDINARY_API_SECRET"),
        },
    }

    log.Println("✅ Config loaded")
}