package repositories

import (
	"context"
	"pbmap_api/src/internal/domain/entities"
	"pbmap_api/src/internal/domain/repositories"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type sessionRepository struct {
	db *gorm.DB
}

func NewSessionRepository(db *gorm.DB) repositories.SessionRepository {
	return &sessionRepository{db: db}
}

func (r *sessionRepository) CreateSession(ctx context.Context, session *entities.UserSession) error {
	return r.db.WithContext(ctx).Create(session).Error
}

func (r *sessionRepository) GetSessionByRefreshToken(ctx context.Context, refreshToken string) (*entities.UserSession, error) {
	var session entities.UserSession
	if err := r.db.WithContext(ctx).Where("refresh_token = ?", refreshToken).First(&session).Error; err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *sessionRepository) GetSessionByDeviceID(ctx context.Context, userID uuid.UUID, deviceID uuid.UUID) (*entities.UserSession, error) {
	var session entities.UserSession
	if err := r.db.WithContext(ctx).Where("user_id = ? AND device_id = ?", userID, deviceID).First(&session).Error; err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *sessionRepository) UpdateSession(ctx context.Context, session *entities.UserSession) error {
	return r.db.WithContext(ctx).Save(session).Error
}

func (r *sessionRepository) RevokeSession(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entities.UserSession{}, id).Error
}
