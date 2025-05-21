package repository

import (
	"database/sql"
	"time"
)

type TokenRepository interface {
	Save(userID uint, token string, expiresAt time.Time) error
	Find(token string) (uint, time.Time, error)
	Delete(token string) error
}

type SQLTokenRepository struct {
	db *sql.DB
}

func NewSQLTokenRepository(db *sql.DB) TokenRepository {
	return &SQLTokenRepository{db: db}
}

func (r *SQLTokenRepository) Save(userID uint, token string, expiresAt time.Time) error {
	_, err := r.db.Exec(
		"INSERT INTO tokens (user_id, token, expires_at) VALUES ($1, $2, $3)",
		userID,
		token,
		expiresAt,
	)
	return err
}

func (r *SQLTokenRepository) Find(token string) (uint, time.Time, error) {
	var (
		userID    uint
		expiresAt time.Time
	)

	err := r.db.QueryRow(
		"SELECT user_id, expires_at FROM tokens WHERE token = $1",
		token,
	).Scan(&userID, &expiresAt)

	return userID, expiresAt, err
}

func (r *SQLTokenRepository) Delete(token string) error {
	_, err := r.db.Exec(
		"DELETE FROM tokens WHERE token = $1",
		token,
	)
	return err
}
