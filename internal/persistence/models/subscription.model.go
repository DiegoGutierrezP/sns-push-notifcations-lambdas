package models

import (
	"fmt"
	"time"
)

const (
	SubscriptionPkPrefix string = "DEVICE#"
	SubscriptionSkPrefix string = "TOPIC#"

	SubscriptionTopicIndex string = "GSI_TopicID"
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

func SubscriptionPk(deviceId string) string {
	return fmt.Sprintf("%s%s", SubscriptionPkPrefix, deviceId)
}

func SubscriptionSk(topic string) string {
	return fmt.Sprintf("%s%s", SubscriptionSkPrefix, topic)
}
