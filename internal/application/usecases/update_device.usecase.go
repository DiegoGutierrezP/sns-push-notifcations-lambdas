package usecases

import (
	"context"
	"fmt"
	"lmbd-digital-push-notifications/internal/application/contracts/repositories"
	"lmbd-digital-push-notifications/internal/application/contracts/services"
	"lmbd-digital-push-notifications/internal/application/dtos"
	appErrors "lmbd-digital-push-notifications/internal/application/errors"
	"lmbd-digital-push-notifications/internal/domain/entities"
	domainErrors "lmbd-digital-push-notifications/internal/domain/errors"
	"lmbd-digital-push-notifications/internal/shared/config"
	"log/slog"
	"net/http"
	"strconv"
	"time"
)

type IUpdateDeviceUseCase interface {
	Execute(ctx context.Context, rq dtos.UpdateDeviceRequest) (*dtos.UpdateDeviceResponse, error)
}

type UpdateDeviceUseCase struct {
	snsService              services.ISnsService
	deviceRepository        repositories.IDeviceRepository
	subscriptionRepository  repositories.ISubscriptionRepository
	optimoveGateway         services.IOptimoveGateway
	optimoveCallbackBaseUrl string
	logger                  *slog.Logger
}

func NewUpdateDeviceUseCase(
	config *config.Config,
	snsService services.ISnsService,
	deviceRepository repositories.IDeviceRepository,
	subscriptionRepository repositories.ISubscriptionRepository,
	optimoveGateway services.IOptimoveGateway,
	logger *slog.Logger,
) IUpdateDeviceUseCase {
	return &UpdateDeviceUseCase{
		snsService:              snsService,
		deviceRepository:        deviceRepository,
		subscriptionRepository:  subscriptionRepository,
		optimoveGateway:         optimoveGateway,
		optimoveCallbackBaseUrl: config.Optimove.CallbackBaseUrl,
		logger:                  logger,
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
		// return nil, fmt.Errorf("An error occurred while searching for the device.")
		return nil, appErrors.NewApplicationError(
			domainErrors.DDBQueryFailed,
			http.StatusBadRequest,
			"An error occurred while searching for the device",
			err,
		)
	}

	if device == nil {
		uc.logger.Error("Device not found",
			"requestId", requestID,
			"deviceId", rq.DeviceId,
		)
		// return nil, fmt.Errorf("Device not found")
		return nil, appErrors.NewApplicationError(
			domainErrors.DEVNotFound,
			http.StatusBadRequest,
			"Device not found",
			err,
		)
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

			// return nil, fmt.Errorf("An Error ocurred while checking if device token exists: %w", err)
			return nil, appErrors.NewApplicationError(
				domainErrors.DDBQueryFailed,
				http.StatusBadRequest,
				"An Error ocurred while checking if device token exists",
				err,
			)
		}
		if exists {
			uc.logger.Warn("token already registered on another device",
				"requestId", requestID,
				"deviceId", device.ID,
				"deviceToken", rq.DeviceToken,
			)

			// return nil, fmt.Errorf("token already registered on another device")
			return nil, appErrors.NewApplicationError(
				domainErrors.DEVTokenAlreadyExists,
				http.StatusBadRequest,
				"token already registered on another device",
				err,
			)
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
		// return nil, fmt.Errorf("An error occurred while updating device entity: %w", err)
		return nil, appErrors.NewApplicationError(
			domainErrors.DDBUpdateFailed,
			http.StatusBadRequest,
			"An error occurred while updating device entity",
			err,
		)
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
			// return nil, fmt.Errorf("An error occurred while updating device endpoint: %w", err)
			return nil, appErrors.NewApplicationError(
				domainErrors.SNSUpdateEndpointFailed,
				http.StatusBadRequest,
				"An error occurred while updating device endpoint",
				err,
			)
		}
	}

	if calimacoIdHasChanged {
		uc.logger.Info("CalimacoId has changed")

		//delete all device subscriptions
		_ = uc.deleteSubscriptions(ctx, device)
	}

	if rq.CalimacoId != nil {
		go uc.optimoveSyncCustomerAttributes(ctx, strconv.Itoa(*device.CalimacoId), device)
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

type attrMapping struct {
	get func(*entities.DeviceEntity) string
}

var optimoveMap = map[string]attrMapping{
	"DEVICE_NAME": {
		get: func(d *entities.DeviceEntity) string { return d.DeviceName },
	},
	"APPLICATION_VERSION": {
		get: func(d *entities.DeviceEntity) string { return d.ApplicationVersion },
	},
	"SYSTEM_VERSION": {
		get: func(d *entities.DeviceEntity) string { return d.SystemVersion },
	},
	"OPERATING_SYSTEM": {
		get: func(d *entities.DeviceEntity) string { return d.OperationSystem },
	},
	"INSTALLATION_DATE": {
		get: func(d *entities.DeviceEntity) string {
			return d.CreatedAt.Format("2006-01-02")
		},
	},
	"INSTALLATION_HOUR": {
		get: func(d *entities.DeviceEntity) string {
			return d.CreatedAt.Format("15:04:05")
		},
	},
}

func (uc *UpdateDeviceUseCase) optimoveSyncCustomerAttributes(ctx context.Context, customerId string, device *entities.DeviceEntity) error {
	uc.logger.Info("optimoveSyncCustomerAttributes started",
		"customerId", customerId,
		"deviceId", device.ID,
	)

	attrs, err := uc.optimoveGateway.GetCustomerAttributes(ctx, customerId)
	if err != nil {

		uc.logger.Error("optimoveGateway.GetCustomerAttributes failed",
			"customerId", customerId,
			"deviceId", device.ID,
			"err", err,
		)

		return err
	}

	var updates []services.OptimoveCustomerAttributeValue

	for _, attr := range attrs {
		mapping, ok := optimoveMap[attr.RealFieldName]
		if !ok {
			continue
		}

		current := fmt.Sprint(attr.Value)
		expected := mapping.get(device)

		if current != expected {
			uc.logger.Info("Optimove attribute out of sync",
				"field", attr.RealFieldName,
				"optimove", current,
				"local", expected,
			)

			updates = append(updates, services.OptimoveCustomerAttributeValue{
				RealFieldName: attr.RealFieldName,
				Value:         expected,
			})
		}
	}

	uc.logger.Info("Total attributes to update",
		"total", len(updates),
	)

	if len(updates) == 0 {
		return nil
	}

	attrValuesUpdate := services.OptimoveCustomerNewAttributesValues{
		CustomerID: customerId,
		Attributes: updates,
	}

	dataUpdate := services.OptimoveUpdateCustomerAttributesDto{
		CustomerNewAttributesValuesList: []services.OptimoveCustomerNewAttributesValues{attrValuesUpdate},
		CallbackURL:                     &uc.optimoveCallbackBaseUrl,
	}

	if err := uc.optimoveGateway.UpdateCustomerAttributes(ctx, dataUpdate); err != nil {
		uc.logger.Error("optimoveGateway.UpdateCustomerAttributes failed",
			"customerId", customerId,
			"deviceId", device.ID,
			"err", err,
		)
		return err
	}

	return nil
}
