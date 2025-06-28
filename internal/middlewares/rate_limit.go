package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

func RateLimitMiddleware(client *redis.Client, limit int, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		key := "rate_limit:" + ip

		current, err := client.Get(c, key).Int()
		if err != nil && err != redis.Nil {
			c.AbortWithStatusJSON(500, gin.H{"error": "Internal server error"})
			return
		}

		if current >= limit {
			c.AbortWithStatusJSON(429, gin.H{"error": "Too many requests"})
			return
		}

		_, err = client.Incr(c, key).Result()
		if err != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": "Internal server error"})
			return
		}

		if current == 0 {
			client.Expire(c, key, window)
		}

		c.Next()
	}
}
