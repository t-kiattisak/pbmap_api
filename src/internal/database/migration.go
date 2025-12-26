package database

import (
	"pbmap_api/src/domain"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&domain.User{},
		&domain.UserSocialAccount{},
		&domain.SpecialCredential{},
		&domain.UserDevice{},
		&domain.UserSession{},
	)
}
