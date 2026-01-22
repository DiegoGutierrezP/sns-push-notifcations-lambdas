package infrastructure

import (
	"lmbd-digital-push-notifications/internal/infrastructure/services/aws"

	"go.uber.org/dig"
)

func RegisterContainer(c *dig.Container) {
	_ = c.Provide(aws.NewSnsService)
}
