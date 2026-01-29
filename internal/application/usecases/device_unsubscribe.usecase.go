package usecases

import (
	"context"
	"errors"
	"lmbd-digital-push-notifications/internal/application/contracts/repositories"
	"lmbd-digital-push-notifications/internal/application/contracts/services"
	"lmbd-digital-push-notifications/internal/application/dtos"
	"log/slog"
)

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
) *DeviceUnsubscribeUseCase {
	return &DeviceUnsubscribeUseCase{
		snsService:             snsService,
		deviceRepository:       deviceRepository,
		subscriptionRepository: subscriptionRepository,
		logger:                 logger,
	}
}

func (uc *DeviceUnsubscribeUseCase) Execute(ctx context.Context, rq dtos.DeviceUnsubscribeRequest) error {
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
		return err
	}

	if len(devices) == 0 {
		uc.logger.Error("No devices found",
			"requestId", requestID,
			"calimacoId", rq.CalimacoId,
		)
		return errors.New("No devices found for calimaco id")
	}

	for _, device := range devices {
		subscription, err := uc.subscriptionRepository.Get(ctx, device.ID.String(), rq.TopicArn)

		if err != nil {
			uc.logger.Error("subscriptionRepository.Get failed:",
				"requestId", requestID,
				"deviceId", device.ID,
				"topicArn", rq.TopicArn,
				"err", err,
			)
			return err
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
			return err
		}

		_ = uc.subscriptionRepository.Delete(ctx, device.ID.String(), rq.TopicArn)

		uc.logger.Info("Unsubscription successfully",
			"requestId", requestID,
			"deviceId", device.ID,
			"subscriptionArn", subscription.SubscriptionArn,
		)
	}

	return nil
}
