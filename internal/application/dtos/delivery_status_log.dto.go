package dtos

import "time"

type SnsDeliveryLog struct {
	Notification struct {
		MessageMD5Sum string `json:"messageMD5Sum"`
		MessageID     string `json:"messageId"`
		TopicArn      string `json:"topicArn"`
		Timestamp     string `json:"timestamp"`
	} `json:"notification"`

	Delivery struct {
		DeliveryID       string `json:"deliveryId"`
		Destination      string `json:"destination"` // endpointArn
		ProviderResponse string `json:"providerResponse"`
		DwellTimeMs      int    `json:"dwellTimeMs"`
		Token            string `json:"token"`
		StatusCode       int    `json:"statusCode"`
	} `json:"delivery"`

	Status string `json:"status"` // SUCCESS / FAILURE
}

type SaveDeliveryStatusInput struct {
	MessageID      string
	EndpointArn    string
	TopicArn       string
	DeliveryStatus string
	StatusCode     int
	Payload        string
	OccurredAt     time.Time
}
