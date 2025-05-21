package router

import (
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/lera-guryan2222/forum/backend/auth-service/internal/controller"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter(authController *controller.AuthController) *gin.Engine {
	r := gin.Default()

	// Настройка CORS с более строгими параметрами
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "X-Requested-With", "Accept"},
		ExposeHeaders:    []string{"Content-Length", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Обработка OPTIONS запросов
	r.OPTIONS("/*any", func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "http://localhost:3000")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization, X-Requested-With, Accept")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Status(200)
	})

	// Логирование всех входящих запросов
	r.Use(func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Log only when path is not healthcheck
		if path != "/health" {
			latency := time.Since(start)
			clientIP := c.ClientIP()
			method := c.Request.Method
			statusCode := c.Writer.Status()
			errorMessage := c.Errors.ByType(gin.ErrorTypePrivate).String()

			if raw != "" {
				path = path + "?" + raw
			}

			log.Printf(
				"[GIN] %v | %3d | %13v | %15s | %-7s %s %s",
				start.Format("2006/01/02 - 15:04:05"),
				statusCode,
				latency,
				clientIP,
				method,
				path,
				errorMessage,
			)
		}
	})

	// Swagger документация
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Группа API для аутентификации
	authGroup := r.Group("/api/auth")
	{
		authGroup.POST("/login", func(c *gin.Context) {
			log.Println("Login request received")
			authController.Login(c)
		})

		authGroup.POST("/register", func(c *gin.Context) {
			log.Println("Register request received with headers:", c.Request.Header)
			authController.Register(c)
		})

		authGroup.POST("/refresh", func(c *gin.Context) {
			log.Println("Refresh request received")
			authController.Refresh(c)
		})
	}

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "Auth service is running",
			"time":    time.Now().Format(time.RFC3339),
		})
	})

	// Обработка 404 ошибок
	r.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{
			"error":   "Not Found",
			"message": "The requested resource was not found",
			"path":    c.Request.URL.Path,
		})
	})

	return r
}
