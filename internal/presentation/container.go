package presentation

import (
	"lmbd-digital-push-notifications/internal/presentation/handlers"

	"go.uber.org/dig"
)

// RegisterContainer registers presentation-layer use cases in the DI container.
func RegisterContainer(c *dig.Container) {
	_ = c.Provide(handlers.NewDeviceSubscribeHandler)
	_ = c.Provide(handlers.NewDeviceUnsubscribeHandler)
	_ = c.Provide(handlers.NewSendNotificationHandler)
	_ = c.Provide(handlers.NewRegisterDeviceHandler)
	_ = c.Provide(handlers.NewUpdateDeviceHandler)
	_ = c.Provide(handlers.NewDeliveryLogProcessorHandler)
}
