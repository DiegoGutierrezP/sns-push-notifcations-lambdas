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

type DeviceUnsubscribeHandler struct {
	usecase usecases.IDeviceUnsubscribeUseCase
}

func NewDeviceUnsubscribeHandler(
	usecase usecases.IDeviceUnsubscribeUseCase,
) *DeviceUnsubscribeHandler {
	return &DeviceUnsubscribeHandler{
		usecase: usecase,
	}
}

func (h *DeviceUnsubscribeHandler) Handler(
	ctx context.Context,
	req events.APIGatewayProxyRequest,
) (events.APIGatewayProxyResponse, error) {
	ctx = context.WithValue(ctx, "RequestID", req.RequestContext.RequestID)

	var payload dtos.DeviceUnsubscribeRequest

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

	resp, err := h.usecase.Execute(ctx, payload)

	if err != nil {
		return httpres.AppFail(err), nil
	}

	var successMessage string
	if resp.Unsubscriptions != nil {
		successMessage = "Dispositivos desuscritos correctamente"
	}

	return httpres.Success(
		http.StatusOK,
		successMessage,
		resp,
	), nil
}
