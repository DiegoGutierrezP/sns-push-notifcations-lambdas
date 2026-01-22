package models

import (
	"time"
)

type SubscriptionModel struct {
	DeviceID        string    `dynamodbav:"pk"`
	TopicID         string    `dynamodbav:"sk"`
	SubscriptionArn string    `dynamodbav:"subscriptionArn"`
	Attributes      string    `dynamodbav:"attributes"`
	IsActive        bool      `dynamodbav:"isActive"`
	CreatedAt       time.Time `dynamodbav:"createdAt"`
	UpdatedAt       time.Time `dynamodbav:"updatedAt"`
}

func (SubscriptionModel) TableName() string {
	return "subscriptions"
}
