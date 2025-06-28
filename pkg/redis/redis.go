package redis

import (
	"context"
	"time"

	"github.com/riad804/go_auth/internal/config"

	"github.com/go-redis/redis/v8"
)

type RedisClient struct {
	*redis.Client
}

func NewRedisClient(cfg *config.Config) (*RedisClient, error) {
	opt, err := redis.ParseURL(cfg.RedisURL)
	if err != nil {
		opt = &redis.Options{
			Addr:     "localhost:6379", // default address
			Password: "",
			DB:       0,
		}
	}

	client := redis.NewClient(opt)
	rc := &RedisClient{Client: client}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if _, err := rc.Ping(ctx).Result(); err != nil {
		return nil, err
	}

	return rc, nil
}

func (rc *RedisClient) Close() error {
	return rc.Client.Close()
}

func (rc *RedisClient) WithContext(ctx context.Context) *RedisClient {
	return &RedisClient{Client: rc.Client.WithContext(ctx)}
}
