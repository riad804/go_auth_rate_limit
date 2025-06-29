package test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	m "github.com/riad804/go_auth/internal/middlewares"
	"github.com/stretchr/testify/assert"
)

func setupTestRedis() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   9, // use a test DB
	})
}

func TestRateLimit(t *testing.T) {
	gin.SetMode(gin.TestMode)
	client := setupTestRedis()
	defer client.FlushDB(client.Context())

	r := gin.New()
	r.POST("/login", m.RateLimitMiddleware(client, 2, time.Minute), func(c *gin.Context) {
		c.JSON(200, gin.H{"msg": "ok"})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/login", nil)

	// should pass
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	// should pass
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	// should be rate limited
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, 429, w.Code)
}
