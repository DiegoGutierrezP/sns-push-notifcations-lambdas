package usecases

import (
	"context"
	"lmbd-digital-push-notifications/internal/application/contracts/repositories"
	"lmbd-digital-push-notifications/internal/application/dtos"
	"log/slog"
)

type SaveDeliveryStatusLogUseCase struct {
	notificationRepository repositories.INotificationRepository
	logger                 *slog.Logger
}

func NewSaveDeliveryStatusLogUseCase(
	logger *slog.Logger,
	notificationRepository repositories.INotificationRepository,
) *SaveDeliveryStatusLogUseCase {
	return &SaveDeliveryStatusLogUseCase{
		logger:                 logger,
		notificationRepository: notificationRepository,
	}
}

func (uc *SaveDeliveryStatusLogUseCase) Execute(ctx context.Context, in dtos.SaveDeliveryStatusInput) error {
	uc.logger.Info("DeliveryStatusLogProcessorUseCase started:",
		"messageId", in.MessageID,
		"endpointArn", in.EndpointArn,
	)

	log := repositories.NotificationLogRegisterDto{
		TopicArn:       in.TopicArn,
		EndpointArn:    in.EndpointArn,
		MessageId:      in.MessageID,
		DeliveryStatus: in.DeliveryStatus,
		StatusCode:     in.StatusCode,
		Payload:        in.Payload,
	}

	if err := uc.notificationRepository.UpsertLog(ctx, &log); err != nil {
		uc.logger.Error("notificationRepository.UpsertLog failed",
			"messageId", in.MessageID,
			"endpointArn", in.EndpointArn,
		)

		return err
	}

	return nil
}
