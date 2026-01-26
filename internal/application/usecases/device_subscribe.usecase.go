package usecases

import (
	"context"
	"lmbd-digital-push-notifications/internal/application/contracts/repositories"
	"lmbd-digital-push-notifications/internal/application/contracts/services"
	"lmbd-digital-push-notifications/internal/application/dtos"
	"log/slog"
)

type DeviceSubscribeUseCase struct {
	snsService       services.ISnsService
	deviceRepository repositories.IDeviceRepository
	logger           *slog.Logger
}

func NewDeviceSubscribeUseCase(
	snsService services.ISnsService,
	deviceRepository repositories.IDeviceRepository,
	logger *slog.Logger,
) *DeviceSubscribeUseCase {
	return &DeviceSubscribeUseCase{
		snsService:       snsService,
		deviceRepository: deviceRepository,
		logger:           logger,
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

	protocol := "application"
	subscriptionAttributes := services.SnsSubscriptionAttributes{
		FilterPolicy: rq.Filters,
	}

	subscriptionArn, err := uc.snsService.Subscription(ctx, rq.TopicArn, device.EndpointArn, &protocol, &subscriptionAttributes)

	if err != nil {
		uc.logger.Error("snsService.Subscription failed",
			"requestId", requestID,
			"err", err,
		)
		return nil, err
	}

	uc.logger.Info("Device subscribed successfully",
		"requestId", requestID,
		"subscriptionArn", subscriptionArn,
	)

	return &dtos.DeviceSubscribeResponse{
		SubscriptionArn: *subscriptionArn,
	}, nil
}
