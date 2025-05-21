package repository

import (
	"database/sql"
	"errors"

	"github.com/lera-guryan2222/forum/backend/auth-service/internal/entity"
)

var ErrRecordNotFound = errors.New("record not found")

type UserRepository interface {
	FindByUsername(username string) (*entity.User, error)
	FindByEmail(email string) (*entity.User, error)
	Create(user *entity.User) error
}

type SQLUserRepository struct {
	db *sql.DB
}

func NewSQLUserRepository(db *sql.DB) UserRepository {
	return &SQLUserRepository{db: db}
}

func (r *SQLUserRepository) FindByUsername(username string) (*entity.User, error) {
	user := &entity.User{}
	err := r.db.QueryRow(
		"SELECT id, username, email, password FROM users WHERE username = $1",
		username,
	).Scan(&user.ID, &user.Username, &user.Email, &user.Password)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *SQLUserRepository) FindByEmail(email string) (*entity.User, error) {
	user := &entity.User{}
	err := r.db.QueryRow(
		"SELECT id, username, email, password FROM users WHERE email = $1",
		email,
	).Scan(&user.ID, &user.Username, &user.Email, &user.Password)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *SQLUserRepository) Create(user *entity.User) error {
	return r.db.QueryRow(
		"INSERT INTO users (username, email, password) VALUES ($1, $2, $3) RETURNING id",
		user.Username,
		user.Email,
		user.Password,
	).Scan(&user.ID)
}
