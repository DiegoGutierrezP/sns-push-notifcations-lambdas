package models

import (
	"time"
)

// NotificationLogModel represents the delivery status log of a message sent to a device
type NotificationLogModel struct {
	PK             string    `dynamodbav:"pk"` // messageId
	SK             string    `dynamodbav:"sk"` // endpointarn
	TopicArn       string    `dynamodbav:"topicArn"`
	DeliveryId     string    `dynamodbav:"deliveryId,omitempty"`
	DeliveryStatus string    `dynamodbav:"deliveryStatus"`
	StatusCode     int       `dynamodbav:"statusCode"`
	Payload        string    `dynamodbav:"payload"`
	CreatedAt      time.Time `dynamodbav:"createdAt"`
	UpdatedAt      time.Time `dynamodbav:"updatedAt"`
}
