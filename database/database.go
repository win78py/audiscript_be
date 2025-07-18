package database

import (
	"audiscript_be/config"
	"context"
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Service interface {
	Health() map[string]string
	DB() *gorm.DB
	Close() error
}

type service struct {
	db *gorm.DB
}

var (
	dbInstance *service
	once       sync.Once
)

func New() Service {
	once.Do(func() {
		dbCfg := config.AppConfig.DB
		dsn := fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s search_path=%s",
			dbCfg.Host,
			dbCfg.User,
			dbCfg.Password,
			dbCfg.Name,
			dbCfg.Port,
			dbCfg.SSLMode,
			dbCfg.Schema,
		)

		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Fatalf("❌ Failed to connect DB: %v", err)
		}

		log.Println("✅ Connected to DB (GORM)")
		dbInstance = &service{db: db}
	})

	return dbInstance
}

func (s *service) DB() *gorm.DB {
	return s.db
}

func (s *service) Close() error {
	sqlDB, err := s.db.DB()
	if err != nil {
		return err
	}
	log.Println("🔌 Closing DB connection...")
	return sqlDB.Close()
}

func (s *service) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	stats := make(map[string]string)

	sqlDB, err := s.db.DB() // Lấy lại sql.DB từ gorm.DB
	if err != nil {
		stats["status"] = "down"
		stats["error"] = fmt.Sprintf("db extract error: %v", err)
		return stats
	}

	if err := sqlDB.PingContext(ctx); err != nil {
		stats["status"] = "down"
		stats["error"] = fmt.Sprintf("db down: %v", err)
		return stats
	}

	dbStats := sqlDB.Stats()
	stats["status"] = "up"
	stats["message"] = "It's healthy"
	stats["open_connections"] = strconv.Itoa(dbStats.OpenConnections)
	stats["in_use"] = strconv.Itoa(dbStats.InUse)
	stats["idle"] = strconv.Itoa(dbStats.Idle)
	stats["wait_count"] = strconv.FormatInt(dbStats.WaitCount, 10)
	stats["wait_duration"] = dbStats.WaitDuration.String()
	stats["max_idle_closed"] = strconv.FormatInt(dbStats.MaxIdleClosed, 10)
	stats["max_lifetime_closed"] = strconv.FormatInt(dbStats.MaxLifetimeClosed, 10)

	return stats
}