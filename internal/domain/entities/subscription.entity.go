package entities

import (
	"time"
)

// SubscriptionEntity represents the subscription of a device to a specific topic.
type SubscriptionEntity struct {
	DeviceId        string
	TopicId         string
	SubscriptionArn string
	Attributes      string
	IsActive        bool
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func NewSubscription(
	deviceId, topicArn, subscriptionArn, attributes string,
) *SubscriptionEntity {

	now := time.Now().UTC()

	return &SubscriptionEntity{
		DeviceId:        deviceId,
		TopicId:         topicArn,
		SubscriptionArn: subscriptionArn,
		Attributes:      attributes,
		IsActive:        true,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
}
