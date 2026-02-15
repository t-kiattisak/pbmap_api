package repositories

import (
	"context"
	"fmt"
	"time"

	"pbmap_api/src/internal/domain/repositories"

	"github.com/redis/go-redis/v9"
)

type tokenRepository struct {
	client *redis.Client
}

func NewTokenRepository(client *redis.Client) repositories.TokenRepository {
	return &tokenRepository{
		client: client,
	}
}

func (r *tokenRepository) SetToken(ctx context.Context, userID string, token string, ttl time.Duration) error {
	key := fmt.Sprintf("user:%s:upstream_token", userID)
	return r.client.Set(ctx, key, token, ttl).Err()
}

func (r *tokenRepository) GetToken(ctx context.Context, userID string) (string, error) {
	key := fmt.Sprintf("user:%s:upstream_token", userID)
	val, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", fmt.Errorf("token not found")
	}
	if err != nil {
		return "", err
	}
	return val, nil
}

func (r *tokenRepository) DeleteToken(ctx context.Context, userID string) error {
	key := fmt.Sprintf("user:%s:upstream_token", userID)
	return r.client.Del(ctx, key).Err()
}

func (r *tokenRepository) SetAppToken(ctx context.Context, userID string, token string, ttl time.Duration) error {
	key := fmt.Sprintf("user:%s:app_token", userID)
	return r.client.Set(ctx, key, token, ttl).Err()
}

func (r *tokenRepository) GetAppToken(ctx context.Context, userID string) (string, error) {
	key := fmt.Sprintf("user:%s:app_token", userID)
	val, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", fmt.Errorf("app token not found")
	}
	if err != nil {
		return "", err
	}
	return val, nil
}

func (r *tokenRepository) DeleteAppToken(ctx context.Context, userID string) error {
	key := fmt.Sprintf("user:%s:app_token", userID)
	return r.client.Del(ctx, key).Err()
}
