package delivery

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lera-guryan2222/forum/backend/forum-service/internal/repository"
	"github.com/lera-guryan2222/forum/backend/forum-service/pkg/auth"
)

type AuthMiddleware struct {
	logger   *log.Logger
	userRepo repository.UserRepository
}

func NewAuthMiddleware(logger *log.Logger, userRepo repository.UserRepository) *AuthMiddleware {
	return &AuthMiddleware{
		logger:   logger,
		userRepo: userRepo,
	}
}

func (m *AuthMiddleware) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization required"})
			return
		}

		claims, err := auth.ValidateToken(tokenString)
		if err != nil {
			m.logger.Printf("Invalid token: %v", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		user, err := m.userRepo.GetByUsername(claims.Username)
		if err != nil {
			m.logger.Printf("User not found: %v", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			return
		}

		ctx := context.WithValue(c.Request.Context(), "userID", user.ID)
		c.Request = c.Request.WithContext(ctx)
		c.Set("userID", user.ID)

		c.Next()
	}
}
