package service

import (
	"github.com/lera-guryan2222/forum/backend/forum-service/internal/repository"
	"github.com/lera-guryan2222/forum/backend/forum-service/pkg/logger"
)

type ForumService interface {
	// TODO: описать методы интерфейса
}

type forumService struct {
	repo   repository.PostRepository
	logger logger.Logger
}

func NewForumService(db interface{}, logger logger.Logger) ForumService {
	// db должен быть *gorm.DB, но для шаблона оставим interface{}
	return &forumService{logger: logger}
}
