package dto

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"

	"pbmap_api/src/internal/domain/entities"
)

type CreatePotentialPointInput struct {
	Name       string          `json:"name" validate:"required"`
	Type       string          `json:"type" validate:"required"`
	Latitude   float64         `json:"latitude" validate:"required"`
	Longitude  float64         `json:"longitude" validate:"required"`
	Properties json.RawMessage `json:"properties"`
}

type UpdatePotentialPointInput struct {
	Name       *string         `json:"name"`
	Type       *string         `json:"type"`
	Latitude   *float64        `json:"latitude"`
	Longitude  *float64        `json:"longitude"`
	Properties json.RawMessage `json:"properties"`
}

type PotentialPointResponse struct {
	ID         uuid.UUID      `json:"id"`
	Name       string         `json:"name"`
	Type       string         `json:"type"`
	Location   Location       `json:"location"`
	Properties datatypes.JSON `json:"properties"`
	CreatorID  *uuid.UUID     `json:"creator_id,omitempty"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
}

type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func ToPotentialPointResponse(pp *entities.PotentialPoint) PotentialPointResponse {
	return PotentialPointResponse{
		ID:   pp.ID,
		Name: pp.Name,
		Type: pp.Type,
		Location: Location{
			Latitude:  pp.Latitude,
			Longitude: pp.Longitude,
		},
		Properties: pp.Properties,
		CreatorID:  &pp.CreatedBy,
		CreatedAt:  pp.CreatedAt,
		UpdatedAt:  pp.UpdatedAt,
	}
}
