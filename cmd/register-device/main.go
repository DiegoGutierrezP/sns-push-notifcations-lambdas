package registerdevice

import (
	"lmbd-digital-push-notifications/internal/application"
	"lmbd-digital-push-notifications/internal/infrastructure"
	"lmbd-digital-push-notifications/internal/persistence"
	"lmbd-digital-push-notifications/internal/presentation"
	"lmbd-digital-push-notifications/internal/presentation/handlers"
	"lmbd-digital-push-notifications/internal/shared/config"

	"github.com/aws/aws-lambda-go/lambda"
	"go.uber.org/dig"
)

var handler *handlers.DeviceSubscribeHandler

func init() {
	c := dig.New()
	c.Provide(config.GetConfig())

	persistence.RegisterContainer(c)
	infrastructure.RegisterContainer(c)
	application.RegisterContainer(c)
	presentation.RegisterContainer(c)

	if err := c.Invoke(func(h *handlers.DeviceSubscribeHandler) {
		handler = h
	}); err != nil {
		panic(err)
	}
}

func main() {
	lambda.Start(handler)
}
