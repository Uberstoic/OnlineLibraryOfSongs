package main

import (
	"fmt"
	"github.com/effectivemobile/music-library/internal/handler"
	"github.com/effectivemobile/music-library/internal/repository"
	"github.com/effectivemobile/music-library/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
)

// @title Music Library API
// @version 1.0
// @description A RESTful API for managing a music library with external API integration
// @host localhost:8080
// @BasePath /api/v1
func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Run database migrations
	migrationURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"))

	m, err := migrate.New("file://migrations", migrationURL)
	if err != nil {
		log.Printf("Migration initialization error: %v", err)
	} else {
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Printf("Migration error: %v", err)
		} else {
			log.Println("Migrations completed successfully")
		}
	}

	// Initialize database connection
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"))

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Initialize repositories and services
	songRepo := repository.NewSongRepository(db)
	musicAPI := service.NewMusicAPIService(os.Getenv("MUSIC_API_URL"))
	songHandler := handler.NewSongHandler(songRepo, musicAPI)

	// Initialize router
	router := gin.Default()

	// API routes
	v1 := router.Group("/api/v1")
	{
		songs := v1.Group("/songs")
		{
			songs.POST("", songHandler.Create)
			songs.GET("", songHandler.List)
			songs.GET("/:id/lyrics", songHandler.GetLyrics)
			songs.PUT("/:id", songHandler.Update)
			songs.DELETE("/:id", songHandler.Delete)
		}
	}

	// Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Start server
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on :%s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
