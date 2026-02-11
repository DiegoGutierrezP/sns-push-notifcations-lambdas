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

type DeviceSubscribeHandler struct {
	usecase usecases.IDeviceSubscribeUseCase
}

func NewDeviceSubscribeHandler(
	usecase usecases.IDeviceSubscribeUseCase,
) *DeviceSubscribeHandler {
	return &DeviceSubscribeHandler{
		usecase: usecase,
	}
}

func (h *DeviceSubscribeHandler) Handler(
	ctx context.Context,
	req events.APIGatewayProxyRequest,
) (events.APIGatewayProxyResponse, error) {
	ctx = context.WithValue(ctx, "RequestID", req.RequestContext.RequestID)

	var payload dtos.DeviceSubscribeRequest

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

	return httpres.Success(
		http.StatusOK,
		"Dispositivos suscritos correctamente",
		data,
	), nil
}
