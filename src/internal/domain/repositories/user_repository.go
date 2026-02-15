package repositories

import (
	"context"
	"pbmap_api/src/internal/domain/entities"

	"github.com/google/uuid"
)

type UserRepository interface {
	Create(ctx context.Context, user *entities.User) error
	FindByID(ctx context.Context, id uuid.UUID) (*entities.User, error)
	Update(ctx context.Context, user *entities.User) error
	Delete(ctx context.Context, id uuid.UUID) error
	FindAll(ctx context.Context) ([]entities.User, error)
	FindBySocialID(ctx context.Context, provider, providerID string) (*entities.User, error)
}
