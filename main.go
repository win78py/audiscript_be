package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"audiscript_be/config"

	"audiscript_be/database"
	"audiscript_be/internal/app"
	"audiscript_be/internal/cloudinary"
	"audiscript_be/internal/middleware"
	"audiscript_be/internal/routes"

	"github.com/gin-gonic/gin"
	// _ "github.com/joho/godotenv/autoload"
	// "audiscript_be/config"
)

func main() {
	fmt.Println("PATH:", os.Getenv("PATH"))
	config.LoadConfig()
	// 1. Load PORT từ env
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dbSvc := database.New()
	defer dbSvc.Close()
	db := dbSvc.DB()

	// Khởi tạo engine “trống”
	r := gin.New()

	// 1) Mặc định Gin không có Logger/Recovery, nên ta phải thêm
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// 2) Thêm CORS custom của mày
	r.Use(middleware.CORSMiddleware())

    r.MaxMultipartMemory = 100 << 20
	// Khởi tạo Cloudinary service
	cldClient, err := cloudinary.NewClient(config.AppConfig.Cloudinary)
	if err != nil {
		log.Fatalf("Failed to create Cloudinary client: %v", err)
	}
	cldSvc := cloudinary.NewService(cldClient)

	deps := &app.AppDependencies{
		DB:         db,
		Cloudinary: cldSvc,
	}
	routes.RegisterAll(r, deps)

	// 4. Tạo HTTP Server
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", port),
		Handler:      r,
		ReadTimeout:  20 * time.Second,
		WriteTimeout: 900 * time.Second,
		IdleTimeout:  time.Minute,
	}

	// 5. Graceful shutdown
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		log.Println("Shutting down server...")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			log.Fatalf("Server forced to shutdown: %v", err)
		}
	}()

	// 6. Run
	log.Printf("Server is running at port %s", port)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Listen: %s\n", err)
	}

	log.Println("Server exiting")
}
