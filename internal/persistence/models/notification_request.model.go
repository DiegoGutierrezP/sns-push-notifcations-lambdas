package models

import (
	"time"
)

// NotificationRequestModel represents a request to send a notification.
type NotificationRequestModel struct {
	PK           string    `dynamodbav:"pk"` // uuid
	MessageId    *string   `dynamodbav:"messageId,omitempty"`
	TopicArn     *string   `dynamodbav:"topicArn,omitempty"`
	TargetArn    *string   `dynamodbav:"targetArn,omitempty"`
	Title        string    `dynamodbav:"title"`
	Body         string    `dynamodbav:"body"`
	Status       string    `dynamodbav:"status"` // PENDING, SENDING, DONE, FAILED
	TotalDevices int       `dynamodbav:"totalDevices"`
	CreatedAt    time.Time `dynamodbav:"createdAt"`
	UpdatedAt    time.Time `dynamodbav:"updatedAt"`
}
