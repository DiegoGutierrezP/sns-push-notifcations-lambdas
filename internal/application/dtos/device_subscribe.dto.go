package dtos

type DeviceSubscribeRequest struct {
	TopicArn   string              `json:"topicArn" validate:"required,min=10"`
	CalimacoId int                 `json:"calimacoId" validate:"required"`
	Filters    map[string][]string `json:"filters" validate:"omitempty,subscription_filters"`
}

type DeviceSubscribeResponse struct {
	CalimacoId    int                   `json:"calimacoId"`
	TopicArn      string                `json:"topicArn"`
	Subscriptions []DeviceSubscribedDto `json:"subscriptions"`
}

type DeviceSubscribedDto struct {
	DeviceId        string `json:"deviceId"`
	SubscriptionArn string `json:"subscriptionArn"`
	Success         bool   `json:"success"`
}
