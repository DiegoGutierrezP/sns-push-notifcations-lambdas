package usecases

import (
	"context"
	"lmbd-digital-push-notifications/internal/application/contracts/repositories"
	"lmbd-digital-push-notifications/internal/application/contracts/services"
	"lmbd-digital-push-notifications/internal/application/dtos"
	appErrors "lmbd-digital-push-notifications/internal/application/errors"
	"lmbd-digital-push-notifications/internal/domain/entities"
	domainErrors "lmbd-digital-push-notifications/internal/domain/errors"
	"log/slog"
	"net/http"
)

type IRegisterDeviceUseCase interface {
	Execute(ctx context.Context, request dtos.RegisterDeviceRequest) (*dtos.RegisterDeviceResponse, error)
}

type RegisterDeviceUseCase struct {
	snsService       services.ISnsService
	deviceRepository repositories.IDeviceRepository
	logger           *slog.Logger
}

func NewRegisterDeviceUseCase(
	snsService services.ISnsService,
	deviceRepository repositories.IDeviceRepository,
	logger *slog.Logger,
) IRegisterDeviceUseCase {
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
		// return nil, fmt.Errorf("An Error occurred while checking if device token exists: %w", err)
		return nil, appErrors.NewApplicationError(
			domainErrors.DDBQueryFailed,
			http.StatusBadRequest,
			"An Error occurred while checking if device token exists",
			err,
		)
	}

	if exists {
		uc.logger.Error("Device token already registered", "requestId", requestID)
		// return nil, fmt.Errorf("token already registered")
		return nil, appErrors.NewApplicationError(
			domainErrors.DEVTokenAlreadyExists,
			http.StatusBadRequest,
			"Device token already registered",
			nil,
		)
	}

	// create endpoint arn for new devices
	endpointArn, err := uc.snsService.CreateEndpoint(ctx, request.DeviceToken)

	if err != nil {
		uc.logger.Error("snsService.createEndpoint failed",
			"requestId", requestID,
			"err", err,
		)
		// return nil, fmt.Errorf("An error occurred while creating endpoint ARN: %w", err)
		return nil, appErrors.NewApplicationError(
			domainErrors.SNSCreateEndpointFailed,
			http.StatusBadRequest,
			"An error occurred while creating endpoint",
			err,
		)
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
		request.CalimacoId,
		nil,
		request.OperatingSystem,
		request.SystemVersion,
	)

	if err := uc.deviceRepository.Save(ctx, deviceEntity); err != nil {
		uc.logger.Error("deviceRepository.Save failed",
			"requestId", requestID,
			"err", err,
		)
		// return nil, fmt.Errorf("An error occurred while registering device: %w", err)
		return nil, appErrors.NewApplicationError(
			domainErrors.DDBInsertFailed,
			http.StatusBadRequest,
			"An error occurred while registering device",
			err,
		)
	}

	uc.logger.Info("Device registered successfully",
		"requestId", requestID,
		"deviceId", deviceEntity.ID,
	)

	return &dtos.RegisterDeviceResponse{
		DeviceId:    deviceEntity.ID.String(),
		EndpointArn: deviceEntity.EndpointArn,
	}, nil
}
