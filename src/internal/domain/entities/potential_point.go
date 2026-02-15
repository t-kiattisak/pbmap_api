package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type PotentialPoint struct {
	ID          uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Name        string         `json:"name"`
	Type        string         `json:"type"`
	Latitude    float64        `json:"latitude"`
	Longitude   float64        `json:"longitude"`
	CreatedYear int            `json:"created_year"`
	CreatedBy   uuid.UUID      `json:"created_by"`
	Properties  datatypes.JSON `json:"properties"` // GORM compatible JSON type

	Creator   *User      `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
	CreatedAt time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt *time.Time `gorm:"index" json:"deleted_at,omitempty"`
}
