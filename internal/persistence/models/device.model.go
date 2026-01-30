package models

import (
	"fmt"
	"time"
)

const (
	DevicePkPrefix     string = "DEVICE#"
	DeviceGSI1PkPrefix string = "TOKEN#"
	DeviceGSI2PkPrefix string = "CALIMACOID#"
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

	GSI1PK string `dynamodbav:"gsi1pk"` // deviceToken
	GSI2PK string `dynamodbav:"gsi2pk"` // calimacId
}

func (DeviceModel) TableName() string {
	return "devices"
}

func DevicePk(deviceId string) string {
	return fmt.Sprintf("%s%s", DevicePkPrefix, deviceId)
}

func DeviceGSI1Pk(token string) string {
	return fmt.Sprintf("%s%s", DeviceGSI1PkPrefix, token)
}

func DeviceGSI2Pk(calimacoId string) string {
	return fmt.Sprintf("%s%s", DeviceGSI2PkPrefix, calimacoId)
}
