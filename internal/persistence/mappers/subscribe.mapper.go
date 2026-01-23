package mappers

import (
	"lmbd-digital-push-notifications/internal/domain/entities"
	"lmbd-digital-push-notifications/internal/persistence/models"
	"strings"
)

func ToSubscribeEntity(model *models.SubscriptionModel) *entities.SubscriptionEntity {
	if model == nil {
		return nil
	}

	return &entities.SubscriptionEntity{
		DeviceId:        strings.TrimPrefix(model.DeviceID, models.SubscriptionPkPrefix),
		TopicId:         strings.TrimPrefix(model.DeviceID, models.SubscriptionSkPrefix),
		SubscriptionArn: model.SubscriptionArn,
		Attributes:      model.Attributes,
		IsActive:        model.IsActive,
		CreatedAt:       model.CreatedAt,
		UpdatedAt:       model.UpdatedAt,
	}
}

func ToSubscribeModel(entity *entities.SubscriptionEntity) *models.SubscriptionModel {
	if entity == nil {
		return nil
	}

	return &models.SubscriptionModel{
		DeviceID:        models.SubscriptionPk(entity.DeviceId),
		TopicID:         models.SubscriptionSk(entity.TopicId),
		SubscriptionArn: entity.SubscriptionArn,
		Attributes:      entity.Attributes,
		IsActive:        entity.IsActive,
		CreatedAt:       entity.CreatedAt,
		UpdatedAt:       entity.UpdatedAt,
	}
}
