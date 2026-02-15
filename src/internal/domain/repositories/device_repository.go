package repositories

import (
	"context"
	"pbmap_api/src/internal/domain/entities"
)

type DeviceRepository interface {
	UpsertDevice(ctx context.Context, device *entities.UserDevice) error
}
