package controller

import (
	"fmt"

	"github.com/lera-guryan2222/forum/backend/forum-service/internal/entity"
	"github.com/lera-guryan2222/forum/backend/forum-service/internal/repository"
)

type PostController interface {
	GetAllPosts() ([]*entity.Post, error)
	GetPostByID(id uint) (*entity.Post, error)
	CreatePost(req *entity.PostRequest, authorID uint) (*entity.Post, error)
	UpdatePost(id uint, req *entity.PostRequest) (*entity.Post, error)
	DeletePost(id uint) error
}

type postController struct {
	repo repository.PostRepository
}

func NewPostController(repo repository.PostRepository) PostController {
	return &postController{repo: repo}
}

func (c *postController) GetAllPosts() ([]*entity.Post, error) {
	return c.repo.GetAll()
}

func (c *postController) GetPostByID(id uint) (*entity.Post, error) {
	return c.repo.GetByID(id) // Используем метод репозитория
}

func (c *postController) CreatePost(req *entity.PostRequest, authorID uint) (*entity.Post, error) {
	post := &entity.Post{
		Title:    req.Title,
		Content:  req.Content,
		AuthorID: authorID,
	}

	if err := c.repo.Create(post); err != nil {
		return nil, err
	}

	return post, nil
}

// post_controller.go
func (c *postController) UpdatePost(id uint, req *entity.PostRequest) (*entity.Post, error) {
	// Убираем лишнее получение поста, так как оно не используется
	updatedPost, err := c.repo.Update(id, req)
	if err != nil {
		return nil, fmt.Errorf("update failed: %w", err)
	}
	return updatedPost, nil
}

func (c *postController) DeletePost(id uint) error {
	return c.repo.Delete(id) // Используем метод репозитория
}
