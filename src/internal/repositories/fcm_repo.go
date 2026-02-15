package repositories

import (
	"context"
	"encoding/json"
	"fmt"

	"pbmap_api/src/internal/domain/entities"
	"pbmap_api/src/internal/domain/repositories"
	"pbmap_api/src/pkg/config"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"google.golang.org/api/option"
)

// Ensure fcmRepo implements repositories.FCMRepository.
var _ repositories.FCMRepository = (*fcmRepo)(nil)

type fcmRepo struct {
	client *messaging.Client
}

// NewFCMRepo creates the FCM repository (implements repositories.FCMRepository).
func NewFCMRepo(cfg *config.Config) (repositories.FCMRepository, error) {
	if cfg.FirebaseCredentialsPath == "" {
		return &fcmRepo{}, nil
	}

	opt := option.WithCredentialsFile(cfg.FirebaseCredentialsPath)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return nil, fmt.Errorf("error initializing firebase app: %v", err)
	}

	client, err := app.Messaging(context.Background())
	if err != nil {
		return nil, fmt.Errorf("error getting messaging client: %v", err)
	}

	return &fcmRepo{client: client}, nil
}

func (s *fcmRepo) BroadcastNotification(ctx context.Context, title, body string) error {
	if s.client == nil {
		return fmt.Errorf("firebase client is not initialized")
	}

	topic := "all_devices"
	message := &messaging.Message{
		Data:         map[string]string{"type": "notification"},
		Notification: &messaging.Notification{Title: title, Body: body},
		Android: &messaging.AndroidConfig{
			Priority: "high",
			Notification: &messaging.AndroidNotification{
				Sound:     "default",
				ChannelID: "high_importance_channel",
			},
		},
		Topic: topic,
	}

	response, err := s.client.Send(ctx, message)
	if err != nil {
		return fmt.Errorf("failed to send message to topic %s: %v", topic, err)
	}
	fmt.Printf("Successfully sent message: %s\n", response)
	return nil
}

func (s *fcmRepo) SendAlarm(ctx context.Context, req *entities.AlarmDispatchRequest) error {
	if s.client == nil {
		return fmt.Errorf("firebase client is not initialized")
	}

	centerJSON, _ := json.Marshal(req.Center)
	message := &messaging.Message{
		Data: map[string]string{
			"type":     "alarm",
			"alarm_id": req.AlarmID,
			"urgency":  req.Urgency,
			"center":   string(centerJSON),
			"signal":   req.Signal,
			"content":  req.Content,
		},
		Android: &messaging.AndroidConfig{
			Priority: "high",
		},
		Topic: "all_devices",
	}

	response, err := s.client.Send(ctx, message)
	if err != nil {
		return fmt.Errorf("failed to send alarm to topic: %v", err)
	}
	fmt.Printf("Successfully sent alarm %s: %s\n", req.AlarmID, response)
	return nil
}

func (s *fcmRepo) SubscribeToTopic(ctx context.Context, tokens []string, topic string) (*entities.TopicManagementResponse, error) {
	if s.client == nil {
		return nil, fmt.Errorf("firebase client is not initialized")
	}

	result := &entities.TopicManagementResponse{
		SuccessTokens: make([]string, 0),
		FailureTokens: make([]entities.TopicManagementError, 0),
	}
	if len(tokens) == 0 {
		return result, nil
	}

	batchSize := 1000
	for i := 0; i < len(tokens); i += batchSize {
		end := i + batchSize
		if end > len(tokens) {
			end = len(tokens)
		}
		batch := tokens[i:end]

		response, err := s.client.SubscribeToTopic(ctx, batch, topic)
		if err != nil {
			return nil, fmt.Errorf("failed to subscribe to topic: %v", err)
		}

		failedIndices := make(map[int]string)
		if response.FailureCount > 0 {
			for _, errWrap := range response.Errors {
				failedIndices[errWrap.Index] = errWrap.Reason
				result.FailureTokens = append(result.FailureTokens, entities.TopicManagementError{
					Token:  batch[errWrap.Index],
					Reason: errWrap.Reason,
				})
			}
		}
		for idx, token := range batch {
			if _, exists := failedIndices[idx]; !exists {
				result.SuccessTokens = append(result.SuccessTokens, token)
			}
		}
	}
	return result, nil
}

func (s *fcmRepo) UnsubscribeFromTopic(ctx context.Context, tokens []string, topic string) (*entities.TopicManagementResponse, error) {
	if s.client == nil {
		return nil, fmt.Errorf("firebase client is not initialized")
	}

	result := &entities.TopicManagementResponse{
		SuccessTokens: make([]string, 0),
		FailureTokens: make([]entities.TopicManagementError, 0),
	}
	if len(tokens) == 0 {
		return result, nil
	}

	batchSize := 1000
	for i := 0; i < len(tokens); i += batchSize {
		end := i + batchSize
		if end > len(tokens) {
			end = len(tokens)
		}
		batch := tokens[i:end]

		response, err := s.client.UnsubscribeFromTopic(ctx, batch, topic)
		if err != nil {
			return nil, fmt.Errorf("failed to unsubscribe from topic: %v", err)
		}

		failedIndices := make(map[int]string)
		if response.FailureCount > 0 {
			for _, errWrap := range response.Errors {
				failedIndices[errWrap.Index] = errWrap.Reason
				result.FailureTokens = append(result.FailureTokens, entities.TopicManagementError{
					Token:  batch[errWrap.Index],
					Reason: errWrap.Reason,
				})
			}
		}
		for idx, token := range batch {
			if _, exists := failedIndices[idx]; !exists {
				result.SuccessTokens = append(result.SuccessTokens, token)
			}
		}
	}
	return result, nil
}
