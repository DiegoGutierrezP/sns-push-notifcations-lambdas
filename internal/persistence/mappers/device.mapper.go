package mappers

import (
	"lmbd-digital-push-notifications/internal/domain/entities"
	"lmbd-digital-push-notifications/internal/persistence/models"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

func ToDeviceEntity(model *models.DeviceModel) *entities.DeviceEntity {
	if model == nil {
		return nil
	}

	calimacoId, _ := strconv.Atoi(*model.CalimacoId)
	userId, _ := strconv.Atoi(*model.UserId)

	return &entities.DeviceEntity{
		ID:                 uuid.MustParse(strings.TrimPrefix(model.ID, models.DevicePkPrefix)),
		DeviceToken:        model.DeviceToken,
		DeviceName:         model.DeviceName,
		EndpointArn:        model.EndpointArn,
		ApplicationVersion: model.ApplicationVersion,
		OperationSystem:    model.OperationSystem,
		SystemVersion:      model.SystemVersion,
		Status:             model.Status,
		CalimacoId:         &calimacoId,
		UserId:             &userId,
		CreatedAt:          model.CreatedAt,
		UpdatedAt:          model.UpdatedAt,
	}
}

func ToDeviceModel(entity *entities.DeviceEntity) *models.DeviceModel {
	if entity == nil {
		return nil
	}

	calimacoIdStr := strconv.Itoa(*entity.CalimacoId)
	userIdStr := strconv.Itoa(*entity.UserId)

	return &models.DeviceModel{
		ID:                 models.DevicePk(entity.ID.String()),
		DeviceToken:        entity.DeviceToken,
		DeviceName:         entity.DeviceName,
		EndpointArn:        entity.EndpointArn,
		ApplicationVersion: entity.ApplicationVersion,
		OperationSystem:    entity.OperationSystem,
		SystemVersion:      entity.SystemVersion,
		Status:             entity.Status,
		CalimacoId:         &calimacoIdStr,
		UserId:             &userIdStr,
		CreatedAt:          entity.CreatedAt,
		UpdatedAt:          entity.UpdatedAt,

		GSI1PK: models.DeviceGSI1Pk(entity.DeviceToken),
		GSI2PK: models.DeviceGSI2Pk(calimacoIdStr),
	}
}
