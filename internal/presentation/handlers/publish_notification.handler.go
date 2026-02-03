package handlers

import (
	"context"
	"encoding/json"
	"lmbd-digital-push-notifications/internal/application/dtos"
	"lmbd-digital-push-notifications/internal/application/usecases"
	validation "lmbd-digital-push-notifications/internal/application/validations"
	"log"

	"github.com/aws/aws-lambda-go/events"
)

type PublishNotificationHandler struct {
	usecase *usecases.PublishNotificationUseCase
}

func NewSendNotificationHandler(
	usecase *usecases.PublishNotificationUseCase,
) *PublishNotificationHandler {
	return &PublishNotificationHandler{
		usecase: usecase,
	}
}

func (h *PublishNotificationHandler) Handler(
	ctx context.Context,
	evt events.SQSEvent,
) {
	for _, rec := range evt.Records {

		log.Println("Processing message:", rec.Body)

		var payload dtos.PublishNotificationRequest

		if err := json.Unmarshal([]byte(rec.Body), &payload); err != nil {
			log.Printf("Request inválido, messageId=%s: %v", rec.MessageId, err)
			continue
		}

		if err := validation.Validate.Struct(payload); err != nil {
			log.Printf("Request inválido validacion fallida, messageId=%s: %v", rec.MessageId, err)
			continue
		}

		if err := h.usecase.Execute(ctx, payload); err != nil {
			log.Printf("falló messageId=%s: %v", rec.MessageId, err)
		}
	}
}
