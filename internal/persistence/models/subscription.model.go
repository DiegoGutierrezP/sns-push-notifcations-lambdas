package models

import (
	"time"
)

const (
	SubscriptionPkPrefix string = "DEVICE#"
	SubscriptionSkPrefix string = "TOPIC#"
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

func SubscriptionPk(id string) string {
	return SubscriptionPkPrefix + id
}

func SubscriptionSk(id string) string {
	return SubscriptionSkPrefix + id
}
