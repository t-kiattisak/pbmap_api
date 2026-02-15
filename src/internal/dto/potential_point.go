package dto

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type CreatePotentialPointInput struct {
	Name        string          `json:"name" validate:"required"`
	Type        string          `json:"type" validate:"required"`
	Latitude    float64         `json:"latitude" validate:"required"`
	Longitude   float64         `json:"longitude" validate:"required"`
	CreatedYear int             `json:"created_year"`
	CreatedBy   string          `json:"created_by"`
	Properties  json.RawMessage `json:"properties"`
}

type UpdatePotentialPointInput struct {
	Name        *string         `json:"name"`
	Type        *string         `json:"type"`
	Latitude    *float64        `json:"latitude"`
	Longitude   *float64        `json:"longitude"`
	CreatedYear *int            `json:"created_year"`
	Properties  json.RawMessage `json:"properties"`
}

type PotentialPointResponse struct {
	ID          uuid.UUID      `json:"id"`
	Name        string         `json:"name"`
	Type        string         `json:"type"`
	Location    Location       `json:"location"`
	CreatedYear int            `json:"created_year"`
	CreatedBy   string         `json:"created_by"`
	Properties  datatypes.JSON `json:"properties"`
	CreatorID   *uuid.UUID     `json:"creator_id,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}
