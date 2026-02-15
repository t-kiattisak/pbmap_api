package repositories

import (
	"context"
	"pbmap_api/src/internal/domain/entities"
)

type FCMRepository interface {
	BroadcastNotification(ctx context.Context, title, body string) error
	SendAlarm(ctx context.Context, req *entities.AlarmDispatchRequest) error
	SubscribeToTopic(ctx context.Context, tokens []string, topic string) (*entities.TopicManagementResponse, error)
	UnsubscribeFromTopic(ctx context.Context, tokens []string, topic string) (*entities.TopicManagementResponse, error)
}
