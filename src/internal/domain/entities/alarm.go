package entities

// AlarmDispatchRequest is the payload for dispatching an alarm (used by delivery + usecase + repo).
type AlarmDispatchRequest struct {
	AlarmID string      `json:"alarm_id"`
	Urgency string      `json:"urgency"`
	Center  AlarmCenter `json:"center"`
	Signal  string      `json:"signal"`
	Content string      `json:"content"`
}

// AlarmCenter is the geographic center for an alarm.
type AlarmCenter struct {
	Lat    float64 `json:"lat"`
	Lng    float64 `json:"lng"`
	Radius int     `json:"radius"`
}

// TopicManagementResponse is the result of subscribe/unsubscribe operations.
type TopicManagementResponse struct {
	SuccessTokens []string               `json:"success_tokens"`
	FailureTokens []TopicManagementError `json:"failure_tokens"`
}

// TopicManagementError represents a failed token in topic operations.
type TopicManagementError struct {
	Token  string `json:"token"`
	Reason string `json:"reason"`
}
