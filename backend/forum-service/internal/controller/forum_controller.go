package controller

import (
	"net/http"

	"github.com/lera-guryan2222/forum/backend/forum-service/internal/service"

	"github.com/gin-gonic/gin"
)

type ForumController struct {
	service service.ForumService
}

func NewForumController(s service.ForumService) *ForumController {
	return &ForumController{service: s}
}

func (c *ForumController) GetPosts(ctx *gin.Context) {
	// TODO: реализовать получение постов
	ctx.JSON(http.StatusOK, gin.H{"message": "get posts endpoint"})
}

func (c *ForumController) CreatePost(ctx *gin.Context) {
	// TODO: реализовать создание поста
	ctx.JSON(http.StatusOK, gin.H{"message": "create post endpoint"})
}
