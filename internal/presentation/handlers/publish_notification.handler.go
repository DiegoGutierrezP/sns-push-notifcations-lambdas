package handlers

import (
	"context"
	"encoding/json"
	"lmbd-digital-push-notifications/internal/application/dtos"
	"lmbd-digital-push-notifications/internal/application/usecases"
	validation "lmbd-digital-push-notifications/internal/application/validations"
	httpres "lmbd-digital-push-notifications/internal/presentation/utils"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

type PublishNotificationHandler struct {
	usecase usecases.IPublishNotificationUseCase
}

func NewSendNotificationHandler(
	usecase usecases.IPublishNotificationUseCase,
) *PublishNotificationHandler {
	return &PublishNotificationHandler{
		usecase: usecase,
	}
}

func (h *PublishNotificationHandler) Handler(
	ctx context.Context,
	req events.APIGatewayProxyRequest,
) (events.APIGatewayProxyResponse, error) {
	ctx = context.WithValue(ctx, "RequestID", req.RequestContext.RequestID)

	var payload dtos.PublishNotificationRequest

	if err := json.Unmarshal([]byte(req.Body), &payload); err != nil {
		return httpres.Fail(http.StatusBadRequest, "Request inválido", err.Error()), nil
	}

	if err := validation.Validate.Struct(payload); err != nil {
		return httpres.Fail(http.StatusUnprocessableEntity, "Validación fallida",
			validation.StringValidationErrors(err),
		), nil
	}

	data, err := h.usecase.Execute(ctx, payload)

	if err != nil {
		log.Printf("error procesando request: %v", err)
		return httpres.Fail(
			http.StatusInternalServerError,
			"error interno",
			err.Error(),
		), nil
	}

	return httpres.Success(http.StatusOK, "", data), nil
}

// func (h *PublishNotificationHandler) Handler(
// 	ctx context.Context,
// 	evt events.SQSEvent,
// ) {
// 	for _, rec := range evt.Records {

// 		log.Println("Processing message:", rec.Body)

// 		var payload dtos.PublishNotificationRequest

// 		if err := json.Unmarshal([]byte(rec.Body), &payload); err != nil {
// 			log.Printf("Request inválido, messageId=%s: %v", rec.MessageId, err)
// 			continue
// 		}

// 		if err := validation.Validate.Struct(payload); err != nil {
// 			log.Printf("Request inválido validacion fallida, messageId=%s: %v", rec.MessageId, err)
// 			continue
// 		}

// 		if err := h.usecase.Execute(ctx, payload); err != nil {
// 			log.Printf("falló messageId=%s: %v", rec.MessageId, err)
// 		}
// 	}
// }
