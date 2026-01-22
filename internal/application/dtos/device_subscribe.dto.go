package dtos

type DeviceSubscribeRequest struct {
	TopicArn    string              `json:"topicArn" validate:"required,min=10"`
	DeviceId    string              `json:"deviceId" validate:"required"`
	DeviceToken string              `json:"deviceToken" validate:"required,min=10"`
	Filters     map[string][]string `json:"filters" validate:"omitempty,subscription_filters"`
}
