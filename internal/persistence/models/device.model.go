package models

import (
	"time"
)

const (
	DeviceCalimacoIdIndex  string = "GSI_CalimacoId"
	DeviceDeviceTokenIndex string = "GSI_DeviceToken"
)

type DeviceModel struct {
	ID                     string    `dynamodbav:"pk"`
	DeviceToken            string    `dynamodbav:"deviceToken"`
	PlatformApplicationArn string    `dynamodbav:"platformApplicationArn"`
	EndpointArn            string    `dynamodbav:"endpointArn"`
	DeviceName             string    `dynamodbav:"deviceName"`
	ApplicationVersion     string    `dynamodbav:"applicationVersion"`
	CalimacoId             *string   `dynamodbav:"calimacoId"`
	UserId                 *string   `dynamodbav:"userId"`
	OperationSystem        string    `dynamodbav:"operatingSystem"`
	SystemVersion          string    `dynamodbav:"systemVersion"`
	Status                 int       `dynamodbav:"status"`
	CreatedAt              time.Time `dynamodbav:"createdAt"`
	UpdatedAt              time.Time `dynamodbav:"updatedAt"`
}
