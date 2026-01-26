package main

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

var handler *handlers.PublishNotificationHandler

func init() {
	c := dig.New()
	c.Provide(config.GetConfig)

	persistence.RegisterContainer(c)
	infrastructure.RegisterContainer(c)
	application.RegisterContainer(c)
	presentation.RegisterContainer(c)

	if err := c.Invoke(func(h *handlers.PublishNotificationHandler) {
		handler = h
	}); err != nil {
		panic(err)
	}
}

func main() {
	lambda.Start(handler)
}

// func main() {
// 	evt := events.SQSEvent{
// 		Records: []events.SQSMessage{
// 			{MessageId: "1", Body: `{"message":"Hola","token":"tok-123"}`},
// 			{MessageId: "2", Body: `{"message":"Hola 2","token":"tok-456"}`},
// 		},
// 	}

// 	handler(context.Background(), evt) // reutiliza tu handler real
// }
