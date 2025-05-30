// internal/repository/user_repository.go
package repository

import (
	"github.com/lera-guryan2222/forum/backend/forum-service/internal/entity"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *entity.User) error
	GetByUsername(username string) (*entity.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *entity.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) GetByUsername(username string) (*entity.User, error) {
	var user entity.User
	if err := r.db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
