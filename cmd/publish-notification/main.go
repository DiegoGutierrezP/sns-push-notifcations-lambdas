package main

import (
	"context"
	"lmbd-digital-push-notifications/internal/application"
	"lmbd-digital-push-notifications/internal/infrastructure"
	"lmbd-digital-push-notifications/internal/persistence"
	"lmbd-digital-push-notifications/internal/presentation"
	"lmbd-digital-push-notifications/internal/presentation/handlers"
	"lmbd-digital-push-notifications/internal/shared/config"

	"github.com/aws/aws-lambda-go/events"
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

// func main() {
// 	lambda.Start(handler.Handler)
// }

func main() {
	evt := events.SQSEvent{
		Records: []events.SQSMessage{
			// {MessageId: "1", Body: `{"targetArn":"arn:aws:sns:us-east-1:418274024107:endpoint/GCM/IncidenciasTest/e7873244-d148-3f70-b4e5-1cb2fc037e61", "body": "Este es el mensaje 1"}`},
			{MessageId: "12", Body: `{"topicArn":"arn:aws:sns:us-east-1:418274024107:football-events", "body": "Este es el mensaje desde un TOPICO"}`},
		},
	}

	handler.Handler(context.Background(), evt) // reutiliza tu handler real
}

// arn:aws:sns:us-east-1:418274024107:endpoint/GCM/IncidenciasTest/e7873244-d148-3f70-b4e5-1cb2fc037e61
// arn:aws:sns:us-east-1:418274024107:endpoint/GCM/IncidenciasTest/cd26bdf3-4526-36a2-a04d-85f2f8de5ac2
