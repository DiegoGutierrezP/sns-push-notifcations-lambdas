package dtos

type PublishNotificationRequest struct {
	TargetArn  string            `json:"targetArn,omitempty"` // endpoint
	TopicArn   string            `json:"topicArn,omitempty"`  // topic
	Title      string            `json:"title"`
	Body       string            `json:"body"`
	Data       map[string]string `json:"data,omitempty"`
	Attributes map[string]string `json:"attributes,omitempty"` // filtros SNS
}
