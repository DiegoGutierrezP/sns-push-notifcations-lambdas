package repositories

import (
	"context"
	"lmbd-digital-push-notifications/internal/domain/entities"
)

type ISubscriptionRepository interface {
	Save(ctx context.Context, s *entities.SubscriptionEntity) error
	Get(ctx context.Context, deviceId, topicId string) (*entities.SubscriptionEntity, error)
	Delete(ctx context.Context, deviceId, topicId string) error
	ListByDevice(ctx context.Context, deviceId string) ([]*entities.SubscriptionEntity, error)
	UpdateStatus(ctx context.Context, deviceId, topicId string, isActive bool) error
	UpdateAttributes(ctx context.Context, deviceId, topicId string, attributes string) error
}
