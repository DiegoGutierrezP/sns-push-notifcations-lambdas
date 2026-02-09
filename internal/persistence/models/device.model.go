package models

import (
	"time"
)

// DeviceModel indexes used to optimize DynamoDB queries.
const (
	DeviceCalimacoIdIndex  string = "GSI_CalimacoId"
	DeviceDeviceTokenIndex string = "GSI_DeviceToken"
)

// DeviceModel represents a device registered to receive push notifications.
// It stores identification, metadata, and DynamoDB attributes associated with
// the device and its relationship to a user.
type DeviceModel struct {
	ID                     string    `dynamodbav:"pk"`
	DeviceToken            string    `dynamodbav:"deviceToken"`
	PlatformApplicationArn string    `dynamodbav:"platformApplicationArn"`
	EndpointArn            string    `dynamodbav:"endpointArn"`
	DeviceName             string    `dynamodbav:"deviceName"`
	ApplicationVersion     string    `dynamodbav:"applicationVersion"`
	CalimacoId             *string   `dynamodbav:"calimacoId,omitempty"`
	UserId                 *string   `dynamodbav:"userId,omitempty"`
	OperationSystem        string    `dynamodbav:"operatingSystem"`
	SystemVersion          string    `dynamodbav:"systemVersion"`
	Status                 int       `dynamodbav:"status"`
	CreatedAt              time.Time `dynamodbav:"createdAt"`
	UpdatedAt              time.Time `dynamodbav:"updatedAt"`
}
