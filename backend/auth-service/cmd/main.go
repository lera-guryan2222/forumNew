package main

import (
	"database/sql"
	"log"
	"os"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/lera-guryan2222/forum/backend/auth-service/internal/controller"
	"github.com/lera-guryan2222/forum/backend/auth-service/internal/repository"
	"github.com/lera-guryan2222/forum/backend/auth-service/internal/router"
	"github.com/lera-guryan2222/forum/backend/auth-service/internal/service"
	"github.com/lera-guryan2222/forum/backend/auth-service/internal/usecase"
	"github.com/lera-guryan2222/forum/backend/auth-service/pkg/auth"
	"github.com/lera-guryan2222/forum/backend/auth-service/pkg/database"
)

func main() {
	dbConfig := database.Config{
		Host:     "localhost", // или ваш хост базы данных
		Port:     "5432",      // стандартный порт PostgreSQL
		User:     "postgres",  // ваш пользователь базы данных
		Password: "postgres",  // ваш пароль
		DBName:   "forum",     // имя вашей базы данных
		SSLMode:  "disable",   // режим SSL
	}
	db, err := database.ConnectWithConfig(dbConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Проверка соединения с базой данных
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	if err := runMigrations(db); err != nil {
		log.Fatalf("Migrations failed: %v", err)
	}

	// Проверка соединения
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	// Инициализация репозитория
	userRepo := repository.NewSQLUserRepository(db)
	tokenRepo := repository.NewSQLTokenRepository(db)

	// Инициализация менеджера токенов
	tokenManager := auth.NewTokenManager(
		os.Getenv("ACCESS_TOKEN_SECRET"),
		os.Getenv("REFRESH_TOKEN_SECRET"),
		24*time.Hour,  // Access token expiry
		720*time.Hour, // Refresh token expiry (30 дней)
	)

	// Инициализация usecase
	authUsecase := usecase.NewAuthUsecase(
		userRepo,
		tokenRepo, // Добавляем tokenRepo
		tokenManager,
	)

	// Инициализация сервиса
	authService := service.NewAuthService(authUsecase)

	// Инициализация контроллера
	authController := controller.NewAuthController(authService)

	// Настройка роутера
	r := router.SetupRouter(authController)

	// Определение порта
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	// Запуск сервера
	log.Printf("Auth Service is running on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

} // Исправленная функция runMigrations
func runMigrations(db *sql.DB) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}

	// Правильный путь к миграциям
	m, err := migrate.NewWithDatabaseInstance(
		"file://C:/Users/usr09/OneDrive/Рабочий%20стол/fooorum/backend/auth-service/internal/migrations", // Исправленный путь
		"postgres",
		driver,
	)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	log.Println("Migrations applied successfully")
	return nil
}
