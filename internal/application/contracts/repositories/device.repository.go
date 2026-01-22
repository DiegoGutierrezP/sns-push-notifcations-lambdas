package repositories

import (
	"context"
	"lmbd-digital-push-notifications/internal/domain/entities"
)

type IDeviceRepository interface {
	Save(ctx context.Context, d *entities.DeviceEntity) error
	GetByID(ctx context.Context, id string) (*entities.DeviceEntity, error)
	Delete(ctx context.Context, id string) error
	UpdateStatus(ctx context.Context, id string, status int) error
	List(ctx context.Context) ([]*entities.DeviceEntity, error)
	ExistsByToken(ctx context.Context, token string) (bool, error)
}
