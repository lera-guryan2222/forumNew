// package database

// import (
// 	"database/sql"
// 	"fmt"
// 	"time"

// 	_ "github.com/lib/pq" // драйвер для PostgreSQL
// )

// type Config struct {
// 	Host     string
// 	Port     string
// 	User     string
// 	Password string
// 	DBName   string
// 	SSLMode  string
// }

// func ConnectWithConfig(config Config) (*sql.DB, error) {
// 	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
// 		config.Host, config.Port, config.User, config.Password, config.DBName, config.SSLMode)

// 	db, err := sql.Open("postgres", connStr)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to open database connection: %v", err)
// 	}

// 	// Установка максимального количества открытых соединений
// 	db.SetMaxOpenConns(25)
// 	// Установка максимального количества простаивающих соединений
// 	db.SetMaxIdleConns(25)
// 	// Установка максимального времени жизни соединения
// 	db.SetConnMaxLifetime(5 * time.Minute)

//		return db, nil
//	}
package database

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq" // драйвер для PostgreSQL
)

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func ConnectWithConfig(config Config) (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.User, config.Password, config.DBName, config.SSLMode)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %v", err)
	}

	// Установка максимального количества открытых соединений
	db.SetMaxOpenConns(25)
	// Установка максимального количества простаивающих соединений
	db.SetMaxIdleConns(25)
	// Установка максимального времени жизни соединения
	db.SetConnMaxLifetime(5 * time.Minute)

	// Проверка соединения
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	return db, nil
}
