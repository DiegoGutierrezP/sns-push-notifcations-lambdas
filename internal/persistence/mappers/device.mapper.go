package mappers

import (
	"lmbd-digital-push-notifications/internal/domain/entities"
	"lmbd-digital-push-notifications/internal/persistence/models"
)

func ToDeviceEntity(model *models.DeviceModel) *entities.DeviceEntity {
	if model == nil {
		return nil
	}

	return &entities.DeviceEntity{
		ID:                 model.ID,
		DeviceToken:        model.DeviceToken,
		DeviceName:         model.DeviceName,
		EndpointArn:        model.EndpointArn,
		ApplicationVersion: model.ApplicationVersion,
		OperationSystem:    model.OperationSystem,
		SystemVersion:      model.SystemVersion,
		Status:             model.Status,
		CalimacoId:         model.CalimacoId,
		UserId:             model.UserId,
		CreatedAt:          model.CreatedAt,
		UpdatedAt:          model.UpdatedAt,
	}
}

func ToDeviceModel(entity *entities.DeviceEntity) *models.DeviceModel {
	if entity == nil {
		return nil
	}

	return &models.DeviceModel{
		ID:                 entity.ID,
		DeviceToken:        entity.DeviceToken,
		DeviceName:         entity.DeviceName,
		EndpointArn:        entity.EndpointArn,
		ApplicationVersion: entity.ApplicationVersion,
		OperationSystem:    entity.OperationSystem,
		SystemVersion:      entity.SystemVersion,
		Status:             entity.Status,
		CalimacoId:         entity.CalimacoId,
		UserId:             entity.UserId,
		CreatedAt:          entity.CreatedAt,
		UpdatedAt:          entity.UpdatedAt,
	}
}
