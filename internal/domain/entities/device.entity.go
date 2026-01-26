package entities

import (
	"time"

	"github.com/google/uuid"
)

type DeviceEntity struct {
	ID                     uuid.UUID
	DeviceToken            string
	PlatformApplicationArn string
	EndpointArn            string
	DeviceName             string
	ApplicationVersion     string
	CalimacoId             string
	UserId                 string
	OperationSystem        string
	SystemVersion          string
	Status                 int
	CreatedAt              time.Time
	UpdatedAt              time.Time
}

func NewDevice(
	deviceToken, endpointArn, deviceName, appVersion,
	calimacoId, userId, os, systemVersion string,
) *DeviceEntity {

	now := time.Now().UTC()

	return &DeviceEntity{
		ID:                 uuid.New(),
		DeviceToken:        deviceToken,
		EndpointArn:        endpointArn,
		DeviceName:         deviceName,
		ApplicationVersion: appVersion,
		CalimacoId:         calimacoId,
		UserId:             userId,
		OperationSystem:    os,
		SystemVersion:      systemVersion,
		Status:             1,
		CreatedAt:          now,
		UpdatedAt:          now,
	}
}

func (e *DeviceEntity) ReadyToUpdate() {

}
