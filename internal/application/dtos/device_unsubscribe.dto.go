package dtos

type DeviceUnsubscribeRequest struct {
	DeviceId string `json:"deviceId" validate:"required,min=10"`
	TopicArn string `json:"topicArn" validate:"required,min=10"`
}
