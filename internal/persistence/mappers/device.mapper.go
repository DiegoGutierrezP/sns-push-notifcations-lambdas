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

	var calimacoId *int
	if model.CalimacoId != nil {
		if v, err := strconv.Atoi(*model.CalimacoId); err == nil {
			calimacoId = &v
		}
	}

	var userId *int
	if model.UserId != nil {
		if v, err := strconv.Atoi(*model.UserId); err == nil {
			userId = &v
		}
	}

	return &entities.DeviceEntity{
		ID:                 uuid.MustParse(strings.TrimPrefix(model.ID, models.DevicePkPrefix)),
		DeviceToken:        model.DeviceToken,
		DeviceName:         model.DeviceName,
		EndpointArn:        model.EndpointArn,
		ApplicationVersion: model.ApplicationVersion,
		OperationSystem:    model.OperationSystem,
		SystemVersion:      model.SystemVersion,
		Status:             model.Status,
		CalimacoId:         calimacoId,
		UserId:             userId,
		CreatedAt:          model.CreatedAt,
		UpdatedAt:          model.UpdatedAt,
	}
}

func ToDeviceModel(entity *entities.DeviceEntity) *models.DeviceModel {
	if entity == nil {
		return nil
	}

	var calimacoIdStr *string
	if entity.CalimacoId != nil {
		v := strconv.Itoa(*entity.CalimacoId)
		calimacoIdStr = &v
	}

	var userIdStr *string
	if entity.UserId != nil {
		v := strconv.Itoa(*entity.UserId)
		userIdStr = &v
	}

	return &models.DeviceModel{
		ID:                 models.DevicePk(entity.ID.String()),
		DeviceToken:        entity.DeviceToken,
		DeviceName:         entity.DeviceName,
		EndpointArn:        entity.EndpointArn,
		ApplicationVersion: entity.ApplicationVersion,
		OperationSystem:    entity.OperationSystem,
		SystemVersion:      entity.SystemVersion,
		Status:             entity.Status,
		CalimacoId:         calimacoIdStr,
		UserId:             userIdStr,
		CreatedAt:          entity.CreatedAt,
		UpdatedAt:          entity.UpdatedAt,

		GSI1PK: models.DeviceGSI1Pk(entity.DeviceToken),
		GSI2PK: models.DeviceGSI2Pk(func() string {
			if calimacoIdStr != nil {
				return *calimacoIdStr
			}
			return ""
		}()),
	}
}
