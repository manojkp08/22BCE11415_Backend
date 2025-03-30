package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/manojkp08/22BCE11415_Backend/internal/database"
	"github.com/redis/go-redis"
	"github.com/yourusername/file-sharing-system/internal/cache"
)

func RateLimit(limit int, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, _ := c.Get("user")
		userID := user.(*database.User).ID
		key := "rate_limit:" + userID

		current, err := cache.RedisClient.Get(cache.ctx, key).Int()
		if err != nil && err != redis.Nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		if current >= limit {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "rate limit exceeded",
			})
			return
		}

		_, err = cache.RedisClient.Incr(cache.ctx, key).Result()
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		if current == 0 {
			cache.RedisClient.Expire(cache.ctx, key, window)
		}

		c.Next()
	}
}
