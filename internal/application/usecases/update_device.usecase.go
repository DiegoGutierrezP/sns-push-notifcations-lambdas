package usecases

import (
	"context"
	"errors"
	"lmbd-digital-push-notifications/internal/application/contracts/repositories"
	"lmbd-digital-push-notifications/internal/application/contracts/services"
	"lmbd-digital-push-notifications/internal/application/dtos"
	"log/slog"
)

type UpdateDeviceUseCase struct {
	snsService       services.ISnsService
	deviceRepository repositories.IDeviceRepository
	logger           *slog.Logger
}

func NewUpdateDeviceUseCase(
	snsService services.ISnsService,
	deviceRepository repositories.IDeviceRepository,
	logger *slog.Logger,
) *UpdateDeviceUseCase {
	return &UpdateDeviceUseCase{
		snsService:       snsService,
		deviceRepository: deviceRepository,
		logger:           logger,
	}
}

func (uc *UpdateDeviceUseCase) Execute(ctx context.Context, rq dtos.UpdateDeviceRequest) (*dtos.UpdateDeviceResponse, error) {
	requestID, _ := ctx.Value("requestID").(string)

	uc.logger.Info("UpdateDeviceUseCase started:",
		"requestId", requestID,
		"deviceName", rq.DeviceName,
		"appVersion", rq.ApplicationVersion,
		"calimacoId", rq.CalimacoId,
		"os", rq.OperatingSystem,
		"osVersion", rq.SystemVersion,
	)

	device, err := uc.deviceRepository.GetByID(ctx, rq.DeviceId)

	if err != nil {
		uc.logger.Error("deviceRepository.ExistsByToken failed",
			"requestId", requestID,
			"err", err,
		)
		return nil, errors.New("An Error ocurred ")
	}

	if rq.DeviceToken != nil {
		exists, err := uc.deviceRepository.ExistsByTokenExcept(ctx, *rq.DeviceToken, device.ID.String())

		if err != nil {
			return nil, errors.New("An Error ocurred ")
		}
		if exists {
			return nil, errors.New("token already registered in another device")
		}
	}

	fieldsToUpdate := repositories.UpdateDeviceFields{
		DeviceName:         rq.DeviceName,
		DeviceToken:        rq.DeviceToken,
		CalimacoId:         rq.CalimacoId,
		ApplicationVersion: rq.ApplicationVersion,
		OperatingSystem:    rq.OperatingSystem,
		SystemVersion:      rq.SystemVersion,
	}

	// update device entity
	if err := uc.deviceRepository.Update(ctx, device.ID.String(), &fieldsToUpdate); err != nil {
		return nil, errors.New("Error al actualizar el device entity")
	}

	// update device token in endpoint
	if rq.DeviceToken != nil && (*rq.DeviceToken != device.DeviceToken) {
		if err := uc.snsService.UpdateEndpoint(ctx, device.EndpointArn, *rq.DeviceToken); err != nil {
			return nil, errors.New("Error al actualizar el device endpoint")
		}
	}

	return &dtos.UpdateDeviceResponse{
		DeviceId: device.ID.String(),
	}, nil
}
