package repository

import (
	"context"
	"pbmap_api/src/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SessionRepository interface {
	CreateSession(ctx context.Context, session *domain.UserSession) error
	GetSessionByRefreshToken(ctx context.Context, refreshToken string) (*domain.UserSession, error)
	GetSessionByDeviceID(ctx context.Context, userID uuid.UUID, deviceID uuid.UUID) (*domain.UserSession, error)
	UpdateSession(ctx context.Context, session *domain.UserSession) error
	RevokeSession(ctx context.Context, id uuid.UUID) error
}

type sessionRepository struct {
	db *gorm.DB
}

func NewSessionRepository(db *gorm.DB) SessionRepository {
	return &sessionRepository{db: db}
}

func (r *sessionRepository) CreateSession(ctx context.Context, session *domain.UserSession) error {
	return r.db.WithContext(ctx).Create(session).Error
}

func (r *sessionRepository) GetSessionByRefreshToken(ctx context.Context, refreshToken string) (*domain.UserSession, error) {
	var session domain.UserSession
	if err := r.db.WithContext(ctx).Where("refresh_token = ?", refreshToken).First(&session).Error; err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *sessionRepository) GetSessionByDeviceID(ctx context.Context, userID uuid.UUID, deviceID uuid.UUID) (*domain.UserSession, error) {
	var session domain.UserSession
	if err := r.db.WithContext(ctx).Where("user_id = ? AND device_id = ?", userID, deviceID).First(&session).Error; err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *sessionRepository) UpdateSession(ctx context.Context, session *domain.UserSession) error {
	return r.db.WithContext(ctx).Save(session).Error
}

func (r *sessionRepository) RevokeSession(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&domain.UserSession{}, id).Error
}
