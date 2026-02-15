package repositories

import (
	"context"
	"pbmap_api/src/internal/domain/entities"

	"github.com/google/uuid"
)

type SessionRepository interface {
	CreateSession(ctx context.Context, session *entities.UserSession) error
	GetSessionByRefreshToken(ctx context.Context, refreshToken string) (*entities.UserSession, error)
	GetSessionByDeviceID(ctx context.Context, userID uuid.UUID, deviceID uuid.UUID) (*entities.UserSession, error)
	UpdateSession(ctx context.Context, session *entities.UserSession) error
	RevokeSession(ctx context.Context, id uuid.UUID) error
}
