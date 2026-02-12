package dtos

type PublishNotificationRequest struct {
	TargetArn  *string           `json:"targetArn,omitempty"` // endpoint
	TopicArn   *string           `json:"topicArn,omitempty"`  // topic
	Title      string            `json:"title" validate:"required"`
	Body       string            `json:"body" validate:"required"`
	Data       map[string]string `json:"data,omitempty"`
	Attributes map[string]string `json:"attributes,omitempty"` // filtros SNS
}

type PushMessage struct {
	Default string `json:"default"`
	GCM     string `json:"GCM"`
	APNS    string `json:"APNS"`
}

type PublishNotificationResponse struct {
	MessageId    string `json:"messageId"`
	TotalDevices *int   `json:"totalDevices"`
}
