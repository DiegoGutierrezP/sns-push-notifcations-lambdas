package repositories

import (
	"context"
	"lmbd-digital-push-notifications/internal/domain/entities"
	"time"
)

type UpdateDeviceFields struct {
	DeviceName         *string
	DeviceToken        *string
	ApplicationVersion *string
	CalimacoId         *int64
	UserId             *string
	OperatingSystem    *string
	SystemVersion      *string
	Status             *int
	UpdatedAt          *time.Time
}

type IDeviceRepository interface {
	Save(ctx context.Context, d *entities.DeviceEntity) error
	GetByID(ctx context.Context, id string) (*entities.DeviceEntity, error)
	Delete(ctx context.Context, id string) error
	UpdateStatus(ctx context.Context, id string, status int) error
	List(ctx context.Context) ([]*entities.DeviceEntity, error)
	ExistsByToken(ctx context.Context, token string) (bool, error)
	ExistsByTokenExcept(ctx context.Context, token string, exceptDeviceID string) (bool, error)
	Update(ctx context.Context, id string, update *UpdateDeviceFields) error
	GetByCalimacoId(ctx context.Context, calimacoId string) ([]entities.DeviceEntity, error)
}
