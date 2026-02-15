package repositories

import (
	"context"
	"pbmap_api/src/internal/domain/entities"

	"github.com/google/uuid"
)

type PotentialPointRepository interface {
	Create(ctx context.Context, pp *entities.PotentialPoint) error
	FindByID(ctx context.Context, id uuid.UUID) (*entities.PotentialPoint, error)
	Update(ctx context.Context, pp *entities.PotentialPoint) error
	Delete(ctx context.Context, id uuid.UUID) error
	FindAll(ctx context.Context) ([]entities.PotentialPoint, error)
}
