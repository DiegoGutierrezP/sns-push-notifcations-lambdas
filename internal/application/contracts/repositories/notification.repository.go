package repositories

import (
	"context"
	"lmbd-digital-push-notifications/internal/domain/constants"
)

type NotificationRequestRegisterDto struct {
	MessageId    *string
	TopicArn     *string
	TargetArn    *string
	Title        string
	Body         string
	Status       constants.NotificationRequestStatus
	TotalDevices int
}

type NotificationRequestUpdateDto struct {
	PK           string
	MessageId    *string
	Status       *constants.NotificationRequestStatus
	TotalDevices *int
}

type NotificationLogRegisterDto struct {
	EndpointArn    string
	TopicArn       string
	DeliveryId     string
	DeliveryStatus string
	MessageId      string
	StatusCode     int
	Payload        string
}

type INotificationRepository interface {
	RegisterRequest(ctx context.Context, nr *NotificationRequestRegisterDto) (*string, error)
	RegisterLog(ctx context.Context, nl *NotificationLogRegisterDto) error
	UpsertLog(ctx context.Context, nl *NotificationLogRegisterDto) error
}
