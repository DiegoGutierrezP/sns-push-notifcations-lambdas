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

type DeviceUnsubscribeHandler struct {
	usecase *usecases.DeviceUnsubscribeUseCase
}

func NewDeviceUnsubscribeHandler(
	usecase *usecases.DeviceUnsubscribeUseCase,
) *DeviceUnsubscribeHandler {
	return &DeviceUnsubscribeHandler{
		usecase: usecase,
	}
}

func (h *DeviceUnsubscribeHandler) Handler(
	ctx context.Context,
	req events.APIGatewayProxyRequest,
) (events.APIGatewayProxyResponse, error) {
	var payload dtos.DeviceUnsubscribeRequest

	if err := json.Unmarshal([]byte(req.Body), &payload); err != nil {
		return httpres.Fail(http.StatusBadRequest, "Request inválido", err.Error()), nil
	}

	if err := validation.Validate.Struct(payload); err != nil {
		return httpres.Fail(http.StatusUnprocessableEntity, "Validación fallida",
			validation.StringValidationErrors(err),
		), nil
	}

	if err := h.usecase.Execute(ctx, payload); err != nil {
		log.Printf("error procesando request: %v", err)
		return httpres.Fail(http.StatusInternalServerError, "error interno", err.Error()), nil
	}

	return httpres.Success[any](http.StatusOK, "", nil), nil
}
