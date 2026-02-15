package entities

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID          uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Email       *string    `gorm:"unique" json:"email"`
	DisplayName string     `json:"display_name"`
	Role        string     `gorm:"type:varchar(20);comment:citizen, officer, admin" json:"role"` // citizen, officer, admin
	CreatedAt   time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt   *time.Time `gorm:"index" json:"deleted_at,omitempty"`

	// Relations
	SocialAccounts    []UserSocialAccount `gorm:"foreignKey:UserID" json:"social_accounts,omitempty"`
	SpecialCredential *SpecialCredential  `gorm:"foreignKey:UserID" json:"special_credential,omitempty"`
	Devices           []UserDevice        `gorm:"foreignKey:UserID" json:"devices,omitempty"`
	Sessions          []UserSession       `gorm:"foreignKey:UserID" json:"sessions,omitempty"`
}

type UserSocialAccount struct {
	ID         uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	UserID     uuid.UUID `json:"user_id"`
	Provider   string    `gorm:"comment:line, google" json:"provider"` // line, google
	ProviderID string    `gorm:"unique;comment:ID from social provider" json:"provider_id"`
}

type SpecialCredential struct {
	UserID       uuid.UUID `gorm:"type:uuid;primaryKey" json:"user_id"`
	Username     string    `gorm:"unique" json:"username"`
	PasswordHash string    `json:"-"`
}

type UserDevice struct {
	ID         uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	UserID     uuid.UUID `json:"user_id"`
	PushToken  string    `gorm:"unique;comment:Generic token for FCM/APNs" json:"push_token"`
	Provider   string    `gorm:"comment:fcm, apns" json:"provider"`            // fcm, apns
	DeviceType string    `gorm:"comment:ios, android, web" json:"device_type"` // ios, android, web
	LastSeen   time.Time `gorm:"default:now()" json:"last_seen"`
}

type UserSession struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	UserID       uuid.UUID `json:"user_id"`
	RefreshToken string    `gorm:"unique" json:"refresh_token"`
	DeviceID     uuid.UUID `json:"device_id"`
	ExpiresAt    time.Time `json:"expires_at"`
	CreatedAt    time.Time `gorm:"default:now()" json:"created_at"`
}
