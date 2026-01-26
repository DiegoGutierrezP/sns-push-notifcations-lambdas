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

type DeviceSubscribeHandler struct {
	usecase *usecases.DeviceSubscribeUseCase
}

func NewDeviceSubscribeHandler(
	usecase *usecases.DeviceSubscribeUseCase,
) *DeviceSubscribeHandler {
	return &DeviceSubscribeHandler{
		usecase: usecase,
	}
}

func (h *DeviceSubscribeHandler) Handler(
	ctx context.Context,
	req events.APIGatewayProxyRequest,
) (events.APIGatewayProxyResponse, error) {

	var payload dtos.DeviceSubscribeRequest

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
		return httpres.Fail(http.StatusInternalServerError, "error interno", err.Error()), nil
	}

	return httpres.Success(http.StatusOK, "", data), nil
}
