package handlers

import (
	"context"
	"encoding/json"
	"lmbd-digital-push-notifications/internal/application/dtos"
	"lmbd-digital-push-notifications/internal/application/usecases"
	validation "lmbd-digital-push-notifications/internal/application/validations"
	"log/slog"

	"github.com/aws/aws-lambda-go/events"
)

type PublishNotificationHandler struct {
	usecase usecases.IPublishNotificationUseCase
	logger  *slog.Logger
}

func NewSendNotificationHandler(
	usecase usecases.IPublishNotificationUseCase,
	logger *slog.Logger,
) *PublishNotificationHandler {
	return &PublishNotificationHandler{
		usecase: usecase,
		logger:  logger,
	}
}

func (h *PublishNotificationHandler) Handler(
	ctx context.Context,
	evt events.SQSEvent,
) {
	for _, rec := range evt.Records {

		h.logger.Info("Processing message",
			"messageId", rec.MessageId,
			"body", rec.Body,
		)

		var payload dtos.PublishNotificationRequest

		if err := json.Unmarshal([]byte(rec.Body), &payload); err != nil {
			h.logger.Error("Invalid request format",
				"messageId", rec.MessageId,
				"err", err,
			)
			continue
		}

		if err := validation.Validate.Struct(payload); err != nil {
			h.logger.Error("Invalid request parameters",
				"messageId", rec.MessageId,
				"validations", validation.StringValidationErrors(err),
				"err", err,
			)
			continue
		}

		_, err := h.usecase.Execute(ctx, payload)

		if err != nil {
			h.logger.Error("usecase excution failed",
				"messageId", rec.MessageId,
				"err", err,
			)
		}
	}
}

// func (h *PublishNotificationHandler) Handler(
// 	ctx context.Context,
// 	req events.APIGatewayProxyRequest,
// ) (events.APIGatewayProxyResponse, error) {
// 	ctx = context.WithValue(ctx, "RequestID", req.RequestContext.RequestID)

// 	var payload dtos.PublishNotificationRequest

// 	if err := json.Unmarshal([]byte(req.Body), &payload); err != nil {
// 		return httpres.AppFail(appErrors.NewApplicationError(
// 			domainErrors.APIInvalidRequestFormat,
// 			http.StatusBadRequest,
// 			"Request inválido",
// 		)), nil
// 	}

// 	if err := validation.Validate.Struct(payload); err != nil {
// 		return httpres.AppFail(appErrors.NewApplicationError(
// 			domainErrors.APIInvalidRequestFormat,
// 			http.StatusUnprocessableEntity,
// 			validation.StringValidationErrors(err),
// 		)), nil
// 	}

// 	data, err := h.usecase.Execute(ctx, payload)

// 	if err != nil {
// 		return httpres.AppFail(err), nil
// 	}

// 	return httpres.Success(http.StatusOK, "Notificación enviada", data), nil
// }
