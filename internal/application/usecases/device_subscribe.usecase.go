package usecases

import (
	"context"
	"encoding/json"
	"lmbd-digital-push-notifications/internal/application/contracts/repositories"
	"lmbd-digital-push-notifications/internal/application/contracts/services"
	"lmbd-digital-push-notifications/internal/application/dtos"
	"lmbd-digital-push-notifications/internal/domain/entities"
	"log/slog"
)

type DeviceSubscribeUseCase struct {
	snsService             services.ISnsService
	deviceRepository       repositories.IDeviceRepository
	subscriptionRepository repositories.ISubscriptionRepository
	logger                 *slog.Logger
}

func NewDeviceSubscribeUseCase(
	snsService services.ISnsService,
	deviceRepository repositories.IDeviceRepository,
	subscriptionRepository repositories.ISubscriptionRepository,
	logger *slog.Logger,
) *DeviceSubscribeUseCase {
	return &DeviceSubscribeUseCase{
		snsService:             snsService,
		deviceRepository:       deviceRepository,
		subscriptionRepository: subscriptionRepository,
		logger:                 logger,
	}
}

func (uc *DeviceSubscribeUseCase) Execute(ctx context.Context, rq dtos.DeviceSubscribeRequest) (*dtos.DeviceSubscribeResponse, error) {
	requestID, _ := ctx.Value("requestID").(string)

	uc.logger.Info("DeviceSubscribeUseCase started:",
		"requestId", requestID,
		"topicArn", rq.TopicArn,
		"deviceId", rq.DeviceId,
	)

	device, err := uc.deviceRepository.GetByID(ctx, rq.DeviceId)

	if err != nil {
		uc.logger.Error("deviceRepository.GetByID failed",
			"requestId", requestID,
			"err", err,
		)
		return nil, err
	}

	subscription, err := uc.subscriptionRepository.Get(ctx, rq.DeviceId, rq.TopicArn)

	subscriptionAttributes := services.SnsSubscriptionAttributes{
		FilterPolicy: rq.Filters,
	}

	attrBytes, err := json.Marshal(subscriptionAttributes)
	if err != nil {
		uc.logger.Error("attributes conversion struct to string failed",
			"requestId", requestID,
			"err", err,
		)
		return nil, err
	}
	subscriptionAttributesString := string(attrBytes)

	var subscriptionArn *string

	if subscription != nil {
		subscriptionArn = &subscription.SubscriptionArn

		uc.logger.Info("Subscription already exists",
			"requestId", requestID,
			"subscriptionArn", subscriptionArn,
			"deviceId", rq.DeviceId,
			"topicArn", rq.TopicArn,
		)

		err := uc.snsService.UpdateSubscriptionAttributes(ctx, *subscriptionArn, &subscriptionAttributes)

		if err != nil {
			uc.logger.Error("snsService.UpdateSubscriptionAttributes failed",
				"requestId", requestID,
				"err", err,
			)
			return nil, err
		}

		if err := uc.subscriptionRepository.UpdateAttributes(ctx, rq.DeviceId, rq.TopicArn, subscriptionAttributesString); err != nil {
			uc.logger.Error("subscriptionRepository.UpdateAttributes failed",
				"requestId", requestID,
				"err", err,
			)
			return nil, err
		}

		uc.logger.Info("Device subscription updated successfully",
			"requestId", requestID,
			"subscriptionArn", subscriptionArn,
		)
	} else {
		subscriptionArn, err = uc.snsService.Subscription(ctx, rq.TopicArn, device.EndpointArn, "application", &subscriptionAttributes)

		if err != nil {
			uc.logger.Error("snsService.Subscription failed",
				"requestId", requestID,
				"err", err,
			)
			return nil, err
		}

		subscriptionEntity := entities.NewSubscription(rq.DeviceId, rq.TopicArn, *subscriptionArn, subscriptionAttributesString)

		if err := uc.subscriptionRepository.Save(ctx, subscriptionEntity); err != nil {
			uc.logger.Error("subscriptionRepository.Save failed",
				"requestId", requestID,
				"err", err,
			)
			return nil, err
		}

		uc.logger.Info("Device subscribed successfully",
			"requestId", requestID,
			"subscriptionArn", subscriptionArn,
		)
	}

	return &dtos.DeviceSubscribeResponse{
		SubscriptionArn: *subscriptionArn,
	}, nil
}
