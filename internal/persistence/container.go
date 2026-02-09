package persistence

import (
	"lmbd-digital-push-notifications/internal/persistence/database"
	"lmbd-digital-push-notifications/internal/persistence/repositories"

	"go.uber.org/dig"
)

// RegisterContainer registers persistence-layer use cases in the DI container.
func RegisterContainer(c *dig.Container) {
	_ = c.Provide(database.NewDynamoDbConnection)
	_ = c.Provide(repositories.NewDeviceRepository)
	_ = c.Provide(repositories.NewSubscriptionRepository)
	_ = c.Provide(repositories.NewNotificationRepository)
}
