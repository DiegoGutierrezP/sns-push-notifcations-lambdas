package usecases

import (
	"context"
	"encoding/json"
	"fmt"
	"lmbd-digital-push-notifications/internal/application/contracts/repositories"
	"lmbd-digital-push-notifications/internal/application/contracts/services"
	"lmbd-digital-push-notifications/internal/application/dtos"
	"lmbd-digital-push-notifications/internal/domain/entities"
	"log/slog"
)

type IDeviceSubscribeUseCase interface {
	Execute(ctx context.Context, rq dtos.DeviceSubscribeRequest) (*dtos.DeviceSubscribeResponse, error)
}

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
) IDeviceSubscribeUseCase {
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
		"calimacoId", rq.CalimacoId,
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

	uc.logger.Info("Total devices found ",
		"requestId", requestID,
		"calimacoId", rq.CalimacoId,
		"totalDevices", len(devices),
	)

	if len(devices) == 0 {
		uc.logger.Error("No devices found",
			"requestId", requestID,
			"calimacoId", rq.CalimacoId,
		)
		return nil, fmt.Errorf("No devices found for calimaco id %d", rq.CalimacoId)
	}

	subscriptionAttributes := services.SnsSubscriptionAttributes{
		FilterPolicy: rq.Filters,
	}

	subscriptions := make([]dtos.DeviceSubscribedDto, 0, len(devices))

	for _, device := range devices {
		deviceSubscribedResult := dtos.DeviceSubscribedDto{
			DeviceId: device.ID.String(),
			Success:  true,
		}

		subscription, err := uc.subscriptionRepository.Get(ctx, device.ID.String(), rq.TopicArn)

		if err != nil {
			uc.logger.Info("subscriptionRepository.Get failed",
				"requestId", requestID,
				"deviceId", device.ID,
				"topicArn", rq.TopicArn,
			)

			deviceSubscribedResult.Success = false
			subscriptions = append(subscriptions, deviceSubscribedResult)
			continue
		}

		// update subscription
		if subscription != nil {
			deviceSubscribedResult.SubscriptionArn = subscription.SubscriptionArn

			uc.logger.Info("Subscription already exists",
				"requestId", requestID,
				"subscriptionArn", subscription.SubscriptionArn,
				"deviceId", device.ID,
				"topicArn", rq.TopicArn,
			)

			err := uc.updateSubscription(ctx, device, rq.TopicArn, subscription.SubscriptionArn, subscriptionAttributes)

			if err != nil {
				deviceSubscribedResult.Success = false
				subscriptions = append(subscriptions, deviceSubscribedResult)
				continue
			}

			uc.logger.Info("Device subscription updated successfully",
				"requestId", requestID,
				"deviceId", device.ID,
				"subscriptionArn", subscription.SubscriptionArn,
			)

			subscriptions = append(subscriptions, deviceSubscribedResult)

			continue
		}

		//register subscription

		subscriptionArn, err := uc.registerSubscription(ctx, device, rq.TopicArn, subscriptionAttributes)

		deviceSubscribedResult.SubscriptionArn = *subscriptionArn

		if err != nil {
			deviceSubscribedResult.Success = false
			subscriptions = append(subscriptions, deviceSubscribedResult)
			continue
		}

		uc.logger.Info("Device subscribed successfully",
			"requestId", requestID,
			"deviceId", device.ID,
			"subscriptionArn", subscriptionArn,
		)

		deviceSubscribedResult.Success = true
		subscriptions = append(subscriptions, deviceSubscribedResult)
	}

	return &dtos.DeviceSubscribeResponse{
		CalimacoId:    rq.CalimacoId,
		TopicArn:      rq.TopicArn,
		Subscriptions: subscriptions,
	}, nil
}

func (uc *DeviceSubscribeUseCase) registerSubscription(
	ctx context.Context,
	device entities.DeviceEntity,
	topicArn string,
	subscriptionAttrs services.SnsSubscriptionAttributes,
) (*string, error) {
	subscriptionArn, err := uc.snsService.Subscription(ctx, topicArn, device.EndpointArn, "application", &subscriptionAttrs)

	if err != nil {
		uc.logger.Error("snsService.Subscription failed",
			"deviceId", device.ID,
			"topicArn", topicArn,
			"err", err,
		)
		return nil, err
	}

	attrBytes, _ := json.Marshal(subscriptionAttrs)

	subscriptionEntity := entities.NewSubscription(device.ID.String(), topicArn, *subscriptionArn, string(attrBytes))

	if err := uc.subscriptionRepository.Save(ctx, subscriptionEntity); err != nil {
		uc.logger.Error("subscriptionRepository.Save failed",
			"deviceId", device.ID,
			"topicArn", topicArn,
			"err", err,
		)
		return nil, err
	}

	return subscriptionArn, nil
}

func (uc *DeviceSubscribeUseCase) updateSubscription(
	ctx context.Context,
	device entities.DeviceEntity,
	topicArn string,
	subscriptionArn string,
	subscriptionAttrs services.SnsSubscriptionAttributes,
) error {
	err := uc.snsService.UpdateSubscriptionAttributes(ctx, subscriptionArn, &subscriptionAttrs)

	if err != nil {
		uc.logger.Error("snsService.UpdateSubscriptionAttributes failed",
			"err", err,
		)
		return err
	}

	attrBytes, _ := json.Marshal(subscriptionAttrs)

	if err := uc.subscriptionRepository.UpdateAttributes(ctx, device.ID.String(), topicArn, string(attrBytes)); err != nil {
		uc.logger.Error("subscriptionRepository.UpdateAttributes failed",
			"err", err,
		)
		return err
	}

	return nil
}
