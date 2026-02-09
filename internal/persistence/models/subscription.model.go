package models

import (
	"fmt"
	"time"
)

const (
	// SubscriptionModel prefixes for partition and sort keys.
	SubscriptionPkPrefix string = "DEVICE#"
	SubscriptionSkPrefix string = "TOPIC#"

	// SubscriptionTopicIndex is a GSI used to optimize queries by topic.
	SubscriptionTopicIndex string = "GSI_TopicID"
)

// SubscriptionModel represents the subscription of a device to a specific topic.
type SubscriptionModel struct {
	DeviceID        string    `dynamodbav:"pk"`
	TopicID         string    `dynamodbav:"sk"`
	SubscriptionArn string    `dynamodbav:"subscriptionArn"`
	Attributes      string    `dynamodbav:"attributes"`
	IsActive        bool      `dynamodbav:"isActive"`
	CreatedAt       time.Time `dynamodbav:"createdAt"`
	UpdatedAt       time.Time `dynamodbav:"updatedAt"`
}

// SubscriptionPk returns the DynamoDB partition key for the given device ID.
func SubscriptionPk(deviceId string) string {
	return fmt.Sprintf("%s%s", SubscriptionPkPrefix, deviceId)
}

// SubscriptionSk returns the DynamoDB sort key for the given topic ID.
func SubscriptionSk(topic string) string {
	return fmt.Sprintf("%s%s", SubscriptionSkPrefix, topic)
}
