package usecase

import (
	"context"
	"fmt"
	"pbmap_api/src/config"
	"pbmap_api/src/internal/dto"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"google.golang.org/api/option"
)

type FCMService interface {
	BroadcastNotification(ctx context.Context, title, body string) error
	SubscribeToTopic(ctx context.Context, tokens []string, topic string) (*dto.TopicManagementResponse, error)
	UnsubscribeFromTopic(ctx context.Context, tokens []string, topic string) (*dto.TopicManagementResponse, error)
}

type fcmService struct {
	client *messaging.Client
}

func NewFCMService(cfg *config.Config) (FCMService, error) {
	if cfg.FirebaseCredentialsPath == "" {
		return &fcmService{}, nil
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

	return &fcmService{
		client: client,
	}, nil
}

func (s *fcmService) BroadcastNotification(ctx context.Context, title, body string) error {
	if s.client == nil {
		return fmt.Errorf("firebase client is not initialized")
	}

	topic := "all_devices"

	message := &messaging.Message{
		Notification: &messaging.Notification{
			Title: title,
			Body:  body,
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

func (s *fcmService) SubscribeToTopic(ctx context.Context, tokens []string, topic string) (*dto.TopicManagementResponse, error) {
	if s.client == nil {
		return nil, fmt.Errorf("firebase client is not initialized")
	}

	result := &dto.TopicManagementResponse{
		SuccessTokens: make([]string, 0),
		FailureTokens: make([]dto.TopicManagementError, 0),
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

				failedToken := batch[errWrap.Index]
				topicError := dto.TopicManagementError{
					Token:  failedToken,
					Reason: errWrap.Reason,
				}
				result.FailureTokens = append(result.FailureTokens, topicError)
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

func (s *fcmService) UnsubscribeFromTopic(ctx context.Context, tokens []string, topic string) (*dto.TopicManagementResponse, error) {
	if s.client == nil {
		return nil, fmt.Errorf("firebase client is not initialized")
	}

	result := &dto.TopicManagementResponse{
		SuccessTokens: make([]string, 0),
		FailureTokens: make([]dto.TopicManagementError, 0),
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

				failedToken := batch[errWrap.Index]
				topicError := dto.TopicManagementError{
					Token:  failedToken,
					Reason: errWrap.Reason,
				}
				result.FailureTokens = append(result.FailureTokens, topicError)
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
