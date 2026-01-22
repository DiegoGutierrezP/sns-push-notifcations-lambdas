package entities

import "time"

type SubscriptionEntity struct {
	DeviceId        string
	TopicId         string
	SubscriptionArn string
	Attributes      string
	IsActive        bool
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
