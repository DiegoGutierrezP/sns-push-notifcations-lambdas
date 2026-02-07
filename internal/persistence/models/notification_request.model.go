package models

import (
	"time"
)

type NotificationRequestModel struct {
	PK           string    `dynamodbav:"pk"` // uuid
	MessageId    *string   `dynamodbav:"messageId"`
	TopicArn     *string   `dynamodbav:"topicArn"`
	TargetArn    *string   `dynamodbav:"targetArn"`
	Title        string    `dynamodbav:"title"`
	Body         string    `dynamodbav:"body"`
	Status       string    `dynamodbav:"status"` // PENDING, SENDING, DONE, FAILED
	TotalDevices int       `dynamodbav:"totalDevices"`
	CreatedAt    time.Time `dynamodbav:"createdAt"`
	UpdatedAt    time.Time `dynamodbav:"updatedAt"`
}
