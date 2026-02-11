package dtos

type DeviceUnsubscribeRequest struct {
	// DeviceId string `json:"deviceId" validate:"required,min=10"`
	CalimacoId int    `json:"calimacoId" validate:"required"`
	TopicArn   string `json:"topicArn" validate:"required,min=10"`
}

type DeviceUnsubscribeResponse struct {
	TopicArn        string                  `json:"topicArn"`
	Unsubscriptions []DeviceUnsubscribedDto `json:"unsubscriptions"`
}

type DeviceUnsubscribedDto struct {
	DeviceId        string `json:"deviceId"`
	SubscriptionArn string `json:"subscriptionArn"`
	Success         bool   `json:"success"`
}
