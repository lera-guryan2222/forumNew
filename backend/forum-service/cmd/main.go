package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/lera-guryan2222/forum/backend/forum-service/internal/controller"
	"github.com/lera-guryan2222/forum/backend/forum-service/internal/delivery"
	"github.com/lera-guryan2222/forum/backend/forum-service/internal/entity"
	"github.com/lera-guryan2222/forum/backend/forum-service/internal/repository"
	"github.com/lera-guryan2222/forum/backend/forum-service/internal/router"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect() (*gorm.DB, error) {
	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "localhost"
	}

	user := os.Getenv("DB_USER")
	if user == "" {
		user = "user"
	}

	password := os.Getenv("DB_PASSWORD")
	if password == "" {
		password = "postgres"
	}

	dbname := os.Getenv("DB_NAME")
	if dbname == "" {
		dbname = "forum"
	}

	port := os.Getenv("DB_PORT")
	if port == "" {
		port = "5432"
	}

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("database connection error: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("database instance error: %w", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("database ping error: %w", err)
	}

	return db, nil
}

func main() {
	logger := log.New(os.Stdout, "[FORUM] ", log.LstdFlags|log.Lshortfile)

	db, err := Connect()
	if err != nil {
		logger.Fatalf("Database connection failed: %v", err)
	}
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	if err := autoMigrate(db); err != nil {
		logger.Fatalf("Migration failed: %v", err)
	}

	// Инициализация репозиториев
	postRepo := repository.NewPostRepository(db)
	userRepo := repository.NewUserRepository(db) // Добавьте реализацию

	// Инициализация контроллеров
	postCtrl := controller.NewPostController(postRepo)
	authCtrl := controller.NewAuthController(userRepo)

	// Middleware
	authMiddleware := delivery.NewAuthMiddleware(logger, userRepo)
	// Роутер
	router := router.SetupRouter(postCtrl, authCtrl, authMiddleware)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	logger.Printf("Server starting on port %s", port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatalf("Server failed: %v", err)
	}
}

func autoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&entity.User{},
		&entity.Post{},
		&entity.ChatMessage{},
		&entity.Token{},
		&entity.EmailVerification{},
	)
}
