package dto

type BroadcastRequest struct {
	Title string `json:"title" validate:"required"`
	Body  string `json:"body" validate:"required"`
}

type SubscribeRequest struct {
	Tokens []string `json:"tokens" validate:"required,min=1"`
}

type TopicManagementError struct {
	Token  string `json:"token"`
	Reason string `json:"reason"`
}

type TopicManagementResponse struct {
	SuccessTokens []string               `json:"success_tokens"`
	FailureTokens []TopicManagementError `json:"failure_tokens"`
}
