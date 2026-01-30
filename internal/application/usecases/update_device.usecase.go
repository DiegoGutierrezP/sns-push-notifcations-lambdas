package usecases

import (
	"context"
	"errors"
	"lmbd-digital-push-notifications/internal/application/contracts/repositories"
	"lmbd-digital-push-notifications/internal/application/contracts/services"
	"lmbd-digital-push-notifications/internal/application/dtos"
	"lmbd-digital-push-notifications/internal/domain/entities"
	"log/slog"
	"time"
)

type UpdateDeviceUseCase struct {
	snsService             services.ISnsService
	deviceRepository       repositories.IDeviceRepository
	subscriptionRepository repositories.ISubscriptionRepository
	logger                 *slog.Logger
}

func NewUpdateDeviceUseCase(
	snsService services.ISnsService,
	deviceRepository repositories.IDeviceRepository,
	subscriptionRepository repositories.ISubscriptionRepository,
	logger *slog.Logger,
) *UpdateDeviceUseCase {
	return &UpdateDeviceUseCase{
		snsService:             snsService,
		deviceRepository:       deviceRepository,
		subscriptionRepository: subscriptionRepository,
		logger:                 logger,
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
		uc.logger.Error("deviceRepository.GetByID failed",
			"requestId", requestID,
			"err", err,
		)
		return nil, errors.New("An error occurred while searching for the device.")
	}

	if device == nil {
		uc.logger.Error("Device not found",
			"requestId", requestID,
			"deviceId", rq.DeviceId,
		)
		return nil, errors.New("Device not found")
	}

	if rq.DeviceToken != nil {
		uc.logger.Info("Update request DeviceToken",
			"requestId", requestID,
			"deviceToken", rq.DeviceToken,
		)

		exists, err := uc.deviceRepository.ExistsByTokenExcept(ctx, *rq.DeviceToken, device.ID.String())

		if err != nil {
			uc.logger.Error("deviceRepository.ExistsByTokenExcept failed",
				"requestId", requestID,
				"deviceId", device.ID,
				"deviceToken", rq.DeviceToken,
			)

			return nil, errors.New("An Error ocurred ")
		}
		if exists {
			uc.logger.Warn("token already registered on another device",
				"requestId", requestID,
				"deviceId", device.ID,
				"deviceToken", rq.DeviceToken,
			)

			return nil, errors.New("token already registered on another device")
		}
	}

	deviceTokenHasChanged := rq.DeviceToken != nil && (*rq.DeviceToken != device.DeviceToken)
	calimacoIdHasChanged := rq.CalimacoId != nil && rq.CalimacoId != device.CalimacoId

	uc.mapDeviceRequestFields(device, rq)

	// update device entity
	if err := uc.deviceRepository.Update(ctx, device.ID.String(), device); err != nil {
		uc.logger.Error("deviceRepository.Update failed",
			"requestId", requestID,
			"deviceId", rq.DeviceId,
		)
		return nil, errors.New("Error al actualizar el device entity")
	}

	// update device token in endpoint
	if deviceTokenHasChanged {

		uc.logger.Info("DeviceToken has changed")

		if err := uc.snsService.UpdateEndpoint(ctx, device.EndpointArn, *rq.DeviceToken); err != nil {
			uc.logger.Error("snsService.UpdateEndpoint failed",
				"requestId", requestID,
				"endpointArn", device.EndpointArn,
				"deviceToken", rq.DeviceToken,
			)
			return nil, errors.New("Error al actualizar el device endpoint")
		}
	}

	if calimacoIdHasChanged {
		uc.logger.Info("CalimacoId has changed")

		//delete all device subscriptions
		_ = uc.deleteSubscriptions(ctx, device)
	}

	return &dtos.UpdateDeviceResponse{
		DeviceId: device.ID.String(),
	}, nil
}

func (uc *UpdateDeviceUseCase) deleteSubscriptions(ctx context.Context, device *entities.DeviceEntity) error {
	deviceId := device.ID.String()
	subscriptions, err := uc.subscriptionRepository.ListByDevice(ctx, deviceId)

	uc.logger.Info("Total current subscriptions ",
		"total", len(subscriptions),
	)

	if err != nil {
		uc.logger.Error("subscriptionRepository.ListByDevice failed",
			"deviceId", deviceId,
			"err", err,
		)

		return err
	}

	for _, sub := range subscriptions {
		if err := uc.snsService.Unsubscription(ctx, sub.SubscriptionArn); err != nil {
			uc.logger.Error("snsService.Unsubscription failed",
				"deviceId", deviceId,
				"subscriptionArn", sub.SubscriptionArn,
				"err", err,
			)
		}
	}

	if err := uc.subscriptionRepository.DeleteByDevice(ctx, deviceId); err != nil {
		uc.logger.Error("subscriptionRepository.DeleteByDevice failed",
			"deviceId", deviceId,
			"err", err,
		)
	}

	return nil
}

func (uc *UpdateDeviceUseCase) mapDeviceRequestFields(device *entities.DeviceEntity, rq dtos.UpdateDeviceRequest) {

	if rq.DeviceName != nil {
		device.DeviceName = *rq.DeviceName
	}
	if rq.DeviceToken != nil {
		device.DeviceToken = *rq.DeviceToken
	}
	if rq.CalimacoId != nil {
		device.CalimacoId = rq.CalimacoId
	}
	if rq.ApplicationVersion != nil {
		device.ApplicationVersion = *rq.ApplicationVersion
	}
	if rq.OperatingSystem != nil {
		device.OperationSystem = *rq.OperatingSystem
	}
	if rq.SystemVersion != nil {
		device.SystemVersion = *rq.SystemVersion
	}

	device.UpdatedAt = time.Now().UTC()
}
