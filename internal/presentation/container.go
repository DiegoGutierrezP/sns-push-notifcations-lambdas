package presentation

import (
	"lmbd-digital-push-notifications/internal/presentation/handlers"

	"go.uber.org/dig"
)

func RegisterContainer(c *dig.Container) {
	_ = c.Provide(handlers.NewDeviceSubscribeHandler)
	_ = c.Provide(handlers.NewDeviceUnsubscribeHandler)
	_ = c.Provide(handlers.NewSendNotificationHandler)
	_ = c.Provide(handlers.NewRegisterDeviceHandler)
}
