package usecase

import (
	"context"

	"pbmap_api/src/internal/domain/entities"
	"pbmap_api/src/internal/domain/repositories"
)

// AlarmUsecase orchestrates alarm dispatch.
type AlarmUsecase interface {
	DispatchAlarm(ctx context.Context, req *entities.AlarmDispatchRequest) error
}

type alarmUsecase struct {
	fcm repositories.FCMRepository
}

// NewAlarmUsecase creates the alarm usecase.
func NewAlarmUsecase(fcm repositories.FCMRepository) AlarmUsecase {
	return &alarmUsecase{fcm: fcm}
}

func (u *alarmUsecase) DispatchAlarm(ctx context.Context, req *entities.AlarmDispatchRequest) error {
	return u.fcm.SendAlarm(ctx, req)
}
