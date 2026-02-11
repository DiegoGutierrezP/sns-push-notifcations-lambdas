package usecases

import (
	"context"
	"lmbd-digital-push-notifications/internal/application/contracts/repositories"
	"lmbd-digital-push-notifications/internal/application/contracts/services"
	"lmbd-digital-push-notifications/internal/application/dtos"
	appErrors "lmbd-digital-push-notifications/internal/application/errors"
	domainErrors "lmbd-digital-push-notifications/internal/domain/errors"
	"log/slog"
	"net/http"
)

type IDeviceUnsubscribeUseCase interface {
	Execute(ctx context.Context, rq dtos.DeviceUnsubscribeRequest) (*dtos.DeviceUnsubscribeResponse, error)
}

type DeviceUnsubscribeUseCase struct {
	snsService             services.ISnsService
	deviceRepository       repositories.IDeviceRepository
	subscriptionRepository repositories.ISubscriptionRepository
	logger                 *slog.Logger
}

func NewDeviceUnsubscribeUseCase(
	snsService services.ISnsService,
	deviceRepository repositories.IDeviceRepository,
	subscriptionRepository repositories.ISubscriptionRepository,
	logger *slog.Logger,
) IDeviceUnsubscribeUseCase {
	return &DeviceUnsubscribeUseCase{
		snsService:             snsService,
		deviceRepository:       deviceRepository,
		subscriptionRepository: subscriptionRepository,
		logger:                 logger,
	}
}

func (uc *DeviceUnsubscribeUseCase) Execute(ctx context.Context, rq dtos.DeviceUnsubscribeRequest) (*dtos.DeviceUnsubscribeResponse, error) {
	requestID, _ := ctx.Value("requestID").(string)

	uc.logger.Info("DeviceUnsubscribeUseCase started:",
		"requestId", requestID,
		"calimacoId", rq.CalimacoId,
		"topicArn", rq.TopicArn,
	)

	devices, err := uc.deviceRepository.GetByCalimacoId(ctx, rq.CalimacoId)

	if err != nil {
		uc.logger.Error("deviceRepository.GetByCalimacoId failed",
			"requestId", requestID,
			"calimacoId", rq.CalimacoId,
			"err", err,
		)
		// return nil, fmt.Errorf("An error occurred, searching devices for calimaco id %d", rq.CalimacoId)
		return nil, appErrors.NewApplicationError(
			domainErrors.DDBQueryFailed,
			http.StatusBadRequest,
			"An error occurred, searching devices",
			err,
		)
	}

	if len(devices) == 0 {
		uc.logger.Error("No devices found",
			"requestId", requestID,
			"calimacoId", rq.CalimacoId,
		)
		// return nil, fmt.Errorf("No devices found for calimaco id %d", rq.CalimacoId)
		return nil, appErrors.NewApplicationError(
			domainErrors.DEVUserWithoutDevices,
			http.StatusNotFound,
			"No devices found for calimaco id",
			nil,
		)
	}

	unsubscriptions := make([]dtos.DeviceUnsubscribedDto, 0, len(devices))

	for _, device := range devices {

		subscription, err := uc.subscriptionRepository.Get(ctx, device.ID.String(), rq.TopicArn)

		if err != nil {
			uc.logger.Error("subscriptionRepository.Get failed:",
				"requestId", requestID,
				"deviceId", device.ID,
				"topicArn", rq.TopicArn,
				"err", err,
			)
			continue
		}

		if subscription == nil {
			uc.logger.Info("Subscription not found",
				"requestId", requestID,
				"deviceId", device.ID,
				"topicArn", rq.TopicArn,
			)
			continue
		}

		deviceUnubscribedResult := dtos.DeviceUnsubscribedDto{
			DeviceId:        device.ID.String(),
			SubscriptionArn: subscription.SubscriptionArn,
			Success:         true,
		}

		uc.logger.Info("SubscriptionArn found",
			"requestId", requestID,
			"subscriptionArn", subscription.SubscriptionArn,
		)

		if err := uc.snsService.Unsubscription(ctx, subscription.SubscriptionArn); err != nil {
			uc.logger.Error("snsService.Unsubscription failed:",
				"requestId", requestID,
				"subscriptionArn", subscription.SubscriptionArn,
				"err", err,
			)

			deviceUnubscribedResult.Success = false
			unsubscriptions = append(unsubscriptions, deviceUnubscribedResult)
			continue
		}

		_ = uc.subscriptionRepository.Delete(ctx, device.ID.String(), rq.TopicArn)

		uc.logger.Info("Unsubscription successfully",
			"requestId", requestID,
			"deviceId", device.ID,
			"subscriptionArn", subscription.SubscriptionArn,
		)

		unsubscriptions = append(unsubscriptions, deviceUnubscribedResult)
	}

	if len(unsubscriptions) == 0 {
		return nil, appErrors.NewApplicationError(
			domainErrors.SUBSNotSubscribed,
			http.StatusNotFound,
			"No subscription found",
			nil,
		)
	}

	return &dtos.DeviceUnsubscribeResponse{
		TopicArn:        rq.TopicArn,
		Unsubscriptions: unsubscriptions,
	}, nil
}
