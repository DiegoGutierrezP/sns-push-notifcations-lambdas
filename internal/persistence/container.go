package persistence

import (
	"lmbd-digital-push-notifications/internal/persistence/database"
	"lmbd-digital-push-notifications/internal/persistence/repositories"

	"go.uber.org/dig"
)

func RegisterContainer(c *dig.Container) {
	_ = c.Provide(database.NewDynamoDbConnection)
	_ = c.Provide(repositories.NewDeviceRepository)
	_ = c.Provide(repositories.NewSubscriptionRepository)
}
