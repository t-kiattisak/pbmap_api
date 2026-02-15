package repositories

import (
	"context"
	"pbmap_api/src/internal/domain/entities"
	"pbmap_api/src/internal/domain/repositories"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type deviceRepository struct {
	db *gorm.DB
}

func NewDeviceRepository(db *gorm.DB) repositories.DeviceRepository {
	return &deviceRepository{db: db}
}

func (r *deviceRepository) UpsertDevice(ctx context.Context, device *entities.UserDevice) error {
	if device.ID != uuid.Nil {
		result := GetDB(ctx, r.db).Model(device).Where("id = ?", device.ID).Updates(device)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected > 0 {
			return nil
		}
	}

	if device.PushToken != "" {
		var existing entities.UserDevice
		if err := GetDB(ctx, r.db).Where("push_token = ?", device.PushToken).First(&existing).Error; err == nil {
			existing.UserID = device.UserID
			existing.LastSeen = time.Now()
			existing.DeviceType = device.DeviceType
			existing.Provider = device.Provider
			device.ID = existing.ID
			return GetDB(ctx, r.db).Save(&existing).Error
		}
	}

	return GetDB(ctx, r.db).Create(device).Error
}
