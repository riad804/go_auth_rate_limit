package repository

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type TokenRepository struct {
	client *redis.Client
}

func NewTokenRepository(client *redis.Client) *TokenRepository {
	return &TokenRepository{client: client}
}

func (r *TokenRepository) StoreRefreshToken(token, userID string, expiry time.Duration) error {
	return r.client.Set(context.Background(), "refresh:"+token, userID, expiry).Err()
}

func (r *TokenRepository) GetRefreshToken(token string) (string, error) {
	return r.client.Get(context.Background(), "refresh:"+token).Result()
}

func (r *TokenRepository) DeleteRefreshToken(token string) error {
	return r.client.Del(context.Background(), "refresh:"+token).Err()
}
