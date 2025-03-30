package middleware

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/manojkp08/22BCE11415_Backend/internal/cache"
	"github.com/manojkp08/22BCE11415_Backend/internal/database"
)

func RateLimit(limit int, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user from context
		user, exists := c.Get("user")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		userID := user.(*database.User).ID
		key := "rate_limit:" + userID

		// Use request context with timeout
		_, cancel := context.WithTimeout(c.Request.Context(), 500*time.Millisecond)
		defer cancel()

		// Get current count
		currentStr, err := cache.Client.Get(key).Result()
		if err != nil && err != redis.Nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}

		current, _ := strconv.Atoi(currentStr)
		if current >= limit {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error":       "rate limit exceeded",
				"retry_after": window.String(),
			})
			return
		}

		// Increment count
		_, err = cache.Client.Incr(key).Result()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}

		// Set expiry if first request
		if current == 0 {
			cache.Client.Expire(key, window)
		}

		c.Next()
	}
}
