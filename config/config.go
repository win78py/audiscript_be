package config

import (
	"fmt"
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


type JWTConfig struct {
    Secret        string
    AccessExpiry  int // phút
    RefreshExpiry int // giờ
}

type Config struct {
    DB         DBConfig
    Cloudinary CloudinaryConfig
    JWT        JWTConfig
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
        JWT: JWTConfig{
            Secret:        os.Getenv("JWT_SECRET"),
            AccessExpiry:  getEnvAsInt("JWT_ACCESS_EXPIRY", 15),
            RefreshExpiry: getEnvAsInt("JWT_REFRESH_EXPIRY", 168),
        },
    }

    log.Println("✅ Config loaded")
}

func getEnvAsInt(name string, defaultVal int) int {
    valStr := os.Getenv(name)
    if valStr == "" {
        return defaultVal
    }
    var val int
    _, err := fmt.Sscanf(valStr, "%d", &val)
    if err != nil {
        return defaultVal
    }
    return val
}