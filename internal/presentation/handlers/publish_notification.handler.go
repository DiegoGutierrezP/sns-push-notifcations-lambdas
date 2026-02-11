package handlers

import (
	"context"
	"encoding/json"
	"lmbd-digital-push-notifications/internal/application/dtos"
	appErrors "lmbd-digital-push-notifications/internal/application/errors"
	"lmbd-digital-push-notifications/internal/application/usecases"
	validation "lmbd-digital-push-notifications/internal/application/validations"
	domainErrors "lmbd-digital-push-notifications/internal/domain/errors"
	httpres "lmbd-digital-push-notifications/internal/presentation/utils"
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
		return httpres.AppFail(appErrors.NewApplicationError(
			domainErrors.APIInvalidRequestFormat,
			http.StatusBadRequest,
			"Request inválido",
		)), nil
	}

	if err := validation.Validate.Struct(payload); err != nil {
		return httpres.AppFail(appErrors.NewApplicationError(
			domainErrors.APIInvalidRequestFormat,
			http.StatusUnprocessableEntity,
			validation.StringValidationErrors(err),
		)), nil
	}

	data, err := h.usecase.Execute(ctx, payload)

	if err != nil {
		return httpres.AppFail(err), nil
	}

	return httpres.Success(http.StatusOK, "Notificación enviada", data), nil
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
