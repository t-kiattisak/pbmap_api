package usecase

import (
	"context"
	"pbmap_api/src/internal/domain/entities"
	"pbmap_api/src/internal/domain/repositories"
	"pbmap_api/src/internal/dto"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type PotentialPointUsecase interface {
	Create(ctx context.Context, input dto.CreatePotentialPointInput, creatorID *uuid.UUID) (*entities.PotentialPoint, error)
	FindByID(ctx context.Context, id uuid.UUID) (*entities.PotentialPoint, error)
	Update(ctx context.Context, id uuid.UUID, input dto.UpdatePotentialPointInput) (*entities.PotentialPoint, error)
	Delete(ctx context.Context, id uuid.UUID) error
	FindAll(ctx context.Context) ([]entities.PotentialPoint, error)
}

type potentialPointUsecase struct {
	repo repositories.PotentialPointRepository
}

func NewPotentialPointUsecase(repo repositories.PotentialPointRepository) PotentialPointUsecase {
	return &potentialPointUsecase{repo: repo}
}

func (u *potentialPointUsecase) Create(ctx context.Context, input dto.CreatePotentialPointInput, creatorID *uuid.UUID) (*entities.PotentialPoint, error) {
	pp := &entities.PotentialPoint{
		Name:        input.Name,
		Type:        input.Type,
		Latitude:    input.Latitude,
		Longitude:   input.Longitude,
		CreatedYear: input.CreatedYear,
		Properties:  datatypes.JSON(input.Properties),
	}

	if creatorID != nil {
		pp.CreatedBy = *creatorID
	} else if input.CreatedBy != "" {
		if id, err := uuid.Parse(input.CreatedBy); err == nil {
			pp.CreatedBy = id
		}
	}

	if err := u.repo.Create(ctx, pp); err != nil {
		return nil, err
	}

	return pp, nil
}

func (u *potentialPointUsecase) FindByID(ctx context.Context, id uuid.UUID) (*entities.PotentialPoint, error) {
	return u.repo.FindByID(ctx, id)
}

func (u *potentialPointUsecase) Update(ctx context.Context, id uuid.UUID, input dto.UpdatePotentialPointInput) (*entities.PotentialPoint, error) {
	pp, err := u.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if input.Name != nil {
		pp.Name = *input.Name
	}
	if input.Type != nil {
		pp.Type = *input.Type
	}
	if input.Latitude != nil {
		pp.Latitude = *input.Latitude
	}
	if input.Longitude != nil {
		pp.Longitude = *input.Longitude
	}
	if input.CreatedYear != nil {
		pp.CreatedYear = *input.CreatedYear
	}
	if input.Properties != nil {
		pp.Properties = datatypes.JSON(input.Properties)
	}

	if err := u.repo.Update(ctx, pp); err != nil {
		return nil, err
	}

	return pp, nil
}

func (u *potentialPointUsecase) Delete(ctx context.Context, id uuid.UUID) error {
	return u.repo.Delete(ctx, id)
}

func (u *potentialPointUsecase) FindAll(ctx context.Context) ([]entities.PotentialPoint, error) {
	return u.repo.FindAll(ctx)
}
