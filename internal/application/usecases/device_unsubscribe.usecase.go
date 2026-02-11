package usecases

import (
	"context"
	"fmt"
	"lmbd-digital-push-notifications/internal/application/contracts/repositories"
	"lmbd-digital-push-notifications/internal/application/contracts/services"
	"lmbd-digital-push-notifications/internal/application/dtos"
	"log/slog"
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
		return nil, fmt.Errorf("An error occurred, searching devices for calimaco id %d", rq.CalimacoId)
	}

	if len(devices) == 0 {
		uc.logger.Error("No devices found",
			"requestId", requestID,
			"calimacoId", rq.CalimacoId,
		)
		return nil, fmt.Errorf("No devices found for calimaco id %d", rq.CalimacoId)
	}

	var devicesUnsubscribed []string

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
			continue
		}

		_ = uc.subscriptionRepository.Delete(ctx, device.ID.String(), rq.TopicArn)

		uc.logger.Info("Unsubscription successfully",
			"requestId", requestID,
			"deviceId", device.ID,
			"subscriptionArn", subscription.SubscriptionArn,
		)

		devicesUnsubscribed = append(devicesUnsubscribed, device.ID.String())
	}

	return &dtos.DeviceUnsubscribeResponse{
		Devices: devicesUnsubscribed,
	}, nil
}
