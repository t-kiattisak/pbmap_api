package repositories

import (
	"context"
	"time"
)

type TokenRepository interface {
	SetToken(ctx context.Context, userID string, token string, ttl time.Duration) error
	GetToken(ctx context.Context, userID string) (string, error)
	DeleteToken(ctx context.Context, userID string) error

	SetAppToken(ctx context.Context, userID string, token string, ttl time.Duration) error
	GetAppToken(ctx context.Context, userID string) (string, error)
	DeleteAppToken(ctx context.Context, userID string) error
}
