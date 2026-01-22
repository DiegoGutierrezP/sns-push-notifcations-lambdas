package models

import (
	"time"

	"github.com/google/uuid"
)

type DeviceModel struct {
	ID                     uuid.UUID `dynamodbav:"pk"`
	DeviceToken            string    `dynamodbav:"deviceToken"`
	PlatformApplicationArn string    `dynamodbav:"platformApplicationArn"`
	EndpointArn            string    `dynamodbav:"endpointArn"`
	DeviceName             string    `dynamodbav:"deviceName"`
	ApplicationVersion     string    `dynamodbav:"applicationVersion"`
	CalimacoId             string    `dynamodbav:"calimacoId"`
	UserId                 string    `dynamodbav:"userId"`
	OperationSystem        string    `dynamodbav:"operatingSystem"`
	SystemVersion          string    `dynamodbav:"systemVersion"`
	Status                 int       `dynamodbav:"status"`
	CreatedAt              time.Time `dynamodbav:"createdAt"`
	UpdatedAt              time.Time `dynamodbav:"updatedAt"`

	GSI1PK string `dynamodbav:"gsi1pk"`
}

func (DeviceModel) TableName() string {
	return "devices"
}
