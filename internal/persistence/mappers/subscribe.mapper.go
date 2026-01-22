package mappers

import (
	"lmbd-digital-push-notifications/internal/domain/entities"
	"lmbd-digital-push-notifications/internal/persistence/models"
)

func ToSubscribeEntity(model *models.SubscriptionModel) *entities.SubscriptionEntity {
	if model == nil {
		return nil
	}

	return &entities.SubscriptionEntity{
		DeviceId:        model.DeviceID,
		TopicId:         model.TopicID,
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
		DeviceID:        entity.DeviceId,
		TopicID:         entity.TopicId,
		SubscriptionArn: entity.SubscriptionArn,
		Attributes:      entity.Attributes,
		IsActive:        entity.IsActive,
		CreatedAt:       entity.CreatedAt,
		UpdatedAt:       entity.UpdatedAt,
	}
}
