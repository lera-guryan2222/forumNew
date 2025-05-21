package router

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/lera-guryan2222/forum/backend/forum-service/internal/controller"
	"github.com/lera-guryan2222/forum/backend/forum-service/internal/delivery"
	"github.com/lera-guryan2222/forum/backend/forum-service/internal/entity"
	"gorm.io/gorm"
)

// SetupRouter создает и настраивает маршруты приложения
func SetupRouter(
	postCtrl controller.PostController,
	authMiddleware *delivery.AuthMiddleware,
) *gin.Engine {
	router := gin.Default()

	// Настройка CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Логирование запросов
	// Логирование запросов
	router.Use(func(c *gin.Context) {
		start := time.Now()
		c.Next()
		latency := time.Since(start)
		log.Printf("%s %s %s %v", c.Request.Method, c.Request.URL.Path, strconv.Itoa(c.Writer.Status()), latency)
	})

	// Группа публичных маршрутов
	public := router.Group("/api/v1")
	{
		public.GET("/posts", getAllPostsHandler(postCtrl))
		public.GET("/posts/:id", getPostByIDHandler(postCtrl))
	}

	// Группа защищенных маршрутов
	protected := router.Group("/api/v1")
	protected.Use(authMiddleware.Handler())
	{
		protected.POST("/posts", createPostHandler(postCtrl))
		protected.PUT("/posts/:id", updatePostHandler(postCtrl))
		protected.DELETE("/posts/:id", deletePostHandler(postCtrl))
	}

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	return router
}

// Обработчики для Gin:

func getAllPostsHandler(ctrl controller.PostController) gin.HandlerFunc {
	return func(c *gin.Context) {
		posts, err := ctrl.GetAllPosts()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "failed to get posts",
				"details": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, posts)
	}
}

func getPostByIDHandler(ctrl controller.PostController) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "invalid post ID",
				"details": err.Error(),
			})
			return
		}

		post, err := ctrl.GetPostByID(uint(id))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "post not found",
				"details": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, post)
	}
}

func createPostHandler(ctrl controller.PostController) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println("[POST] Attempting to create new post")

		// Логируем входящий запрос
		log.Printf("Incoming request: %s %s", c.Request.Method, c.Request.URL.Path)
		log.Printf("Headers: %+v", c.Request.Header)

		var req entity.PostRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			log.Printf("[ERROR] Invalid request body: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "invalid request body",
				"details": err.Error(),
			})
			return
		}

		log.Printf("Parsed request: Title='%s', Content='%s'", req.Title, req.Content)

		// Получаем ID автора из контекста
		authUserID, exists := c.Get("userID")
		if !exists {
			log.Println("[ERROR] User ID not found in context")
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "user not authenticated",
			})
			return
		}

		authorID, ok := authUserID.(uint)
		if !ok {
			log.Printf("[ERROR] Invalid user ID type: %T", authUserID)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "invalid user ID format",
			})
			return
		}

		log.Printf("Attempting to create post for user ID: %d", authorID)

		resp, err := ctrl.CreatePost(&req, authorID)
		if err != nil {
			log.Printf("[ERROR] Failed to create post: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "failed to create post",
				"details": err.Error(),
			})
			return
		}

		log.Printf("Successfully created post with ID: %d", resp.ID)
		c.JSON(http.StatusCreated, resp)
	}
}

// router.go
func updatePostHandler(ctrl controller.PostController) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid post ID format"})
			return
		}

		var req entity.PostRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "validation error",
				"details": err.Error(),
			})
			return
		}

		resp, err := ctrl.UpdatePost(uint(id), &req)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "update failed",
				"details": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, resp)
	}
}

func deletePostHandler(ctrl controller.PostController) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "invalid post ID",
				"details": err.Error(),
			})
			return
		}

		if err := ctrl.DeletePost(uint(id)); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "failed to delete post",
				"details": err.Error(),
			})
			return
		}
		c.Status(http.StatusNoContent)
	}
}
