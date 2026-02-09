package application

import (
	"lmbd-digital-push-notifications/internal/application/usecases"

	"go.uber.org/dig"
)

// RegisterContainer registers application-layer use cases in the DI container.
func RegisterContainer(c *dig.Container) {
	_ = c.Provide(usecases.NewPublishNotificationUseCase)
	_ = c.Provide(usecases.NewDeviceSubscribeUseCase)
	_ = c.Provide(usecases.NewDeviceUnsubscribeUseCase)
	_ = c.Provide(usecases.NewRegisterDeviceUseCase)
	_ = c.Provide(usecases.NewUpdateDeviceUseCase)
	_ = c.Provide(usecases.NewSaveDeliveryStatusLogUseCase)
}
