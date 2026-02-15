package database

import (
	"pbmap_api/src/internal/domain/entities"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&entities.User{},
		&entities.UserSocialAccount{},
		&entities.SpecialCredential{},
		&entities.UserDevice{},
		&entities.UserSession{},
		&entities.PotentialPoint{},
	)
}
