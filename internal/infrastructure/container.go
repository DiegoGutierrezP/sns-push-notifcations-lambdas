package infrastructure

import (
	"lmbd-digital-push-notifications/internal/infrastructure/logs"
	"lmbd-digital-push-notifications/internal/infrastructure/rates"
	"lmbd-digital-push-notifications/internal/infrastructure/services/aws"
	"lmbd-digital-push-notifications/internal/infrastructure/services/gateways"

	"go.uber.org/dig"
)

func RegisterContainer(c *dig.Container) {
	_ = c.Provide(aws.NewSnsService)
	_ = c.Provide(logs.NewLogger)
	_ = c.Provide(rates.NewRateLimiter)
	_ = c.Provide(gateways.NewOptimoveGateway)
}
