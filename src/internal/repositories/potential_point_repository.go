package repositories

import (
	"context"
	"pbmap_api/src/internal/domain/entities"
	"pbmap_api/src/internal/domain/repositories"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type potentialPointRepository struct {
	db *gorm.DB
}

func NewPotentialPointRepository(db *gorm.DB) repositories.PotentialPointRepository {
	return &potentialPointRepository{db: db}
}

func (r *potentialPointRepository) Create(ctx context.Context, pp *entities.PotentialPoint) error {
	return r.db.WithContext(ctx).Create(pp).Error
}

func (r *potentialPointRepository) FindByID(ctx context.Context, id uuid.UUID) (*entities.PotentialPoint, error) {
	var pp entities.PotentialPoint
	if err := r.db.WithContext(ctx).Preload("Creator").First(&pp, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &pp, nil
}

func (r *potentialPointRepository) Update(ctx context.Context, pp *entities.PotentialPoint) error {
	return r.db.WithContext(ctx).Save(pp).Error
}

func (r *potentialPointRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entities.PotentialPoint{}, "id = ?", id).Error
}

func (r *potentialPointRepository) FindAll(ctx context.Context) ([]entities.PotentialPoint, error) {
	var pps []entities.PotentialPoint
	if err := r.db.WithContext(ctx).Preload("Creator").Find(&pps).Error; err != nil {
		return nil, err
	}
	return pps, nil
}
