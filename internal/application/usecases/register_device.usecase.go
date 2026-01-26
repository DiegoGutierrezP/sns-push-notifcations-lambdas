package usecases

import (
	"context"
	"errors"
	"lmbd-digital-push-notifications/internal/application/contracts/repositories"
	"lmbd-digital-push-notifications/internal/application/contracts/services"
	"lmbd-digital-push-notifications/internal/application/dtos"
	"lmbd-digital-push-notifications/internal/domain/entities"
	"log/slog"
	"strconv"
)

type RegisterDeviceUseCase struct {
	snsService       services.ISnsService
	deviceRepository repositories.IDeviceRepository
	logger           *slog.Logger
}

func NewRegisterDeviceUseCase(
	snsService services.ISnsService,
	deviceRepository repositories.IDeviceRepository,
	logger *slog.Logger,
) *RegisterDeviceUseCase {
	return &RegisterDeviceUseCase{
		snsService:       snsService,
		deviceRepository: deviceRepository,
		logger:           logger,
	}
}

func (uc *RegisterDeviceUseCase) Execute(ctx context.Context, request dtos.RegisterDeviceRequest) (*dtos.RegisterDeviceResponse, error) {
	requestID, _ := ctx.Value("requestID").(string)

	uc.logger.Info("RegisterDeviceUseCase started:",
		"requestId", requestID,
		"deviceName", request.DeviceName,
		"appVersion", request.ApplicationVersion,
		"calimacoId", request.CalimacoId,
		"os", request.OperatingSystem,
		"osVersion", request.SystemVersion,
	)

	exists, err := uc.deviceRepository.ExistsByToken(ctx, request.DeviceToken)

	if err != nil {
		uc.logger.Error("deviceRepository.ExistsByToken failed",
			"requestId", requestID,
			"err", err,
		)
		return nil, errors.New("An Error ocurred ")
	}

	if exists {
		uc.logger.Error("Device token already registered", "requestId", requestID)
		return nil, errors.New("token already registered")
	}

	// create endpoint arn for new devices
	endpointArn, err := uc.snsService.CreateEndpoint(ctx, request.DeviceToken)

	if err != nil {
		uc.logger.Error("snsService.createEndpoint failed",
			"requestId", requestID,
			"err", err,
		)
		return nil, errors.New("Ocurrio un error al crear el endpoint ARN")
	}

	uc.logger.Info("Endpoint created",
		"requestId", requestID,
		"endpointArn", endpointArn,
	)

	// create new device instance
	deviceEntity := entities.NewDevice(
		request.DeviceToken,
		endpointArn,
		request.DeviceName,
		request.ApplicationVersion,
		strconv.Itoa(int(request.CalimacoId)),
		"",
		request.OperatingSystem,
		request.SystemVersion,
	)

	if err := uc.deviceRepository.Save(ctx, deviceEntity); err != nil {
		uc.logger.Error("deviceRepository.Save failed",
			"requestId", requestID,
			"err", err,
		)
		return nil, errors.New("Ocurrio un error al registrar el dispositivo")
	}

	uc.logger.Info("Device registered successfully",
		"requestId", requestID,
		"deviceId", deviceEntity.ID,
	)

	return &dtos.RegisterDeviceResponse{
		DeviceId:    deviceEntity.ID.String(),
		DeviceToken: deviceEntity.DeviceToken,
		EndpointArn: deviceEntity.EndpointArn,
	}, nil
}
