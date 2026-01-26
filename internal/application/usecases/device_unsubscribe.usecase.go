package usecases

import (
	"context"
	"lmbd-digital-push-notifications/internal/application/contracts/repositories"
	"lmbd-digital-push-notifications/internal/application/contracts/services"
	"lmbd-digital-push-notifications/internal/application/dtos"
	"log/slog"
)

type DeviceUnsubscribeUseCase struct {
	snsService             services.ISnsService
	subscriptionRepository repositories.ISubscriptionRepository
	logger                 *slog.Logger
}

func NewDeviceUnsubscribeUseCase(
	snsService services.ISnsService,
	subscriptionRepository repositories.ISubscriptionRepository,
	logger *slog.Logger,
) *DeviceUnsubscribeUseCase {
	return &DeviceUnsubscribeUseCase{
		snsService:             snsService,
		subscriptionRepository: subscriptionRepository,
		logger:                 logger,
	}
}

func (uc *DeviceUnsubscribeUseCase) Execute(ctx context.Context, rq dtos.DeviceUnsubscribeRequest) error {
	requestID, _ := ctx.Value("requestID").(string)

	uc.logger.Info("DeviceUnsubscribeUseCase started:",
		"requestId", requestID,
		"deviceId", rq.DeviceId,
		"topicArn", rq.TopicArn,
	)

	subscription, err := uc.subscriptionRepository.Get(ctx, rq.DeviceId, rq.TopicArn)

	if err != nil {
		uc.logger.Error("subscriptionRepository.Get failed:",
			"requestId", requestID,
			"deviceId", rq.DeviceId,
			"topicArn", rq.TopicArn,
			"err", err,
		)
		return err
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

	uc.logger.Info("Unsubscription successfully",
		"requestId", requestID,
		"subscriptionArn", subscription.SubscriptionArn,
	)

	return nil
}
