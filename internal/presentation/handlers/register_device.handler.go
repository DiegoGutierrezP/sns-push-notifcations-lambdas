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

type RegisterDeviceHandler struct {
	usecase *usecases.RegisterDeviceUseCase
}

func NewRegisterDeviceHandler(
	usecase *usecases.RegisterDeviceUseCase,
) *RegisterDeviceHandler {
	return &RegisterDeviceHandler{
		usecase: usecase,
	}
}

func (h *RegisterDeviceHandler) Handler(
	ctx context.Context,
	req events.APIGatewayProxyRequest,
) (events.APIGatewayProxyResponse, error) {
	ctx = context.WithValue(ctx, "RequestID", req.RequestContext.RequestID)

	var payload dtos.RegisterDeviceRequest

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
