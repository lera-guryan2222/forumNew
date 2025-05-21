package repository

import (
	"github.com/lera-guryan2222/forum/backend/forum-service/internal/entity"
	"gorm.io/gorm"
)

type PostRepository interface {
	Create(post *entity.Post) error
	GetAll() ([]*entity.Post, error)
	GetByID(id uint) (*entity.Post, error) // Добавляем новые методы
	Update(id uint, req *entity.PostRequest) (*entity.Post, error)
	Delete(id uint) error
}

type postRepository struct {
	db *gorm.DB
}

func (r *postRepository) GetByID(id uint) (*entity.Post, error) {
	var post entity.Post
	if err := r.db.Preload("Author").First(&post, id).Error; err != nil {
		return nil, err
	}
	return &post, nil
}

func (r *postRepository) Update(id uint, req *entity.PostRequest) (*entity.Post, error) {
	var post entity.Post
	if err := r.db.First(&post, id).Error; err != nil {
		return nil, err
	}
	updates := map[string]interface{}{
		"Title":   req.Title,
		"Content": req.Content,
	}

	if err := r.db.Model(&post).Updates(updates).Error; err != nil {
		return nil, err
	}

	return &post, nil
}

func (r *postRepository) Delete(id uint) error {
	return r.db.Delete(&entity.Post{}, id).Error
}

func NewPostRepository(db *gorm.DB) PostRepository {
	return &postRepository{db: db}
}

func (r *postRepository) Create(post *entity.Post) error {
	return r.db.Create(post).Error
}

func (r *postRepository) GetAll() ([]*entity.Post, error) {
	var posts []*entity.Post
	err := r.db.Preload("Author").Find(&posts).Error
	return posts, err
}
func (r *postRepository) GetAllWithPagination(offset, limit int) ([]*entity.Post, error) {
	var posts []*entity.Post
	err := r.db.Preload("Author").
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&posts).Error
	return posts, err
}
