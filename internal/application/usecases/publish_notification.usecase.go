package usecases

import (
	"context"
	"fmt"
	"lmbd-digital-push-notifications/internal/application/contracts/repositories"
	"lmbd-digital-push-notifications/internal/application/contracts/services"
	"lmbd-digital-push-notifications/internal/application/dtos"
	"lmbd-digital-push-notifications/internal/domain/constants"
	"log/slog"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns/types"
)

type PublishNotificationUseCase struct {
	snsService             services.ISnsService
	notificationRepository repositories.INotificationRepository
	subscriptionRepository repositories.ISubscriptionRepository
	logger                 *slog.Logger
}

func NewPublishNotificationUseCase(
	snsService services.ISnsService,
	notificationRepository repositories.INotificationRepository,
	subscriptionRepository repositories.ISubscriptionRepository,
	logger *slog.Logger,
) *PublishNotificationUseCase {
	return &PublishNotificationUseCase{
		snsService:             snsService,
		notificationRepository: notificationRepository,
		subscriptionRepository: subscriptionRepository,
		logger:                 logger,
	}
}

func (uc *PublishNotificationUseCase) Execute(ctx context.Context, rq dtos.PublishNotificationRequest) error {
	uc.logger.Info("PublishNotificationUseCase started:",
		"TopicArn", rq.TopicArn,
		"TargetArn", rq.TargetArn,
	)

	if rq.TargetArn == nil && rq.TopicArn == nil {
		uc.logger.Error("targetArn and topicArn are empty")

		return fmt.Errorf("targetArn or topicArn is required")
	}

	pushMessage, publishOptions := uc.buildPushMessage(rq.Title, rq.Body, rq.Data, rq.Attributes)

	notificationRequestDto := repositories.NotificationRequestRegisterDto{
		TopicArn:  rq.TopicArn,
		TargetArn: rq.TargetArn,
		Title:     *rq.Title,
		Body:      rq.Body,
		Status:    constants.NotificationRequestStatusPending,
	}

	messageId, err := uc.snsService.Publish(ctx, rq.TopicArn, rq.TargetArn, pushMessage, &publishOptions)

	if err != nil {
		uc.logger.Error("snsService.Publish failed",
			"TopicArn", rq.TopicArn,
			"TargetArn", rq.TargetArn,
			"err", err,
		)

		notificationRequestDto.Status = constants.NotificationRequestStatusFailed

		return fmt.Errorf("An error occurred while publishing notification: %w", err)
	}

	if rq.TopicArn != nil {
		countSubs, err := uc.subscriptionRepository.CountByTopic(ctx, *rq.TopicArn)

		if err != nil {
			uc.logger.Error("subscriptionRepository.CountByTopic failed",
				"TopicArn", rq.TopicArn,
				"err", err,
			)
		} else {
			notificationRequestDto.TotalDevices = int(countSubs)
		}
	}

	notificationRequestDto.MessageId = messageId
	notificationRequestDto.Status = constants.NotificationRequestStatusSent

	_, err = uc.notificationRepository.RegisterRequest(ctx, &notificationRequestDto)

	if err != nil {
		uc.logger.Error("notificationRepository.RegisterRequest failed",
			"topicArn", rq.TopicArn,
			"targetArn", rq.TargetArn,
			"err", err,
		)
	}

	uc.logger.Info("Notification published successfully",
		"topicArn", rq.TopicArn,
		"targetArn", rq.TargetArn,
	)

	return nil
}

func (uc *PublishNotificationUseCase) buildPushMessage(
	title *string,
	body string,
	data map[string]string,
	attributes map[string]string,
) (dtos.PushMessage, services.SnsPublishOptions) {

	// Android (Firebase / GCM)
	gcm := map[string]any{
		"notification": map[string]any{
			"title": title,
			"body":  body,
		},
		"data": data,
	}

	// iOS (APNS)
	apns := map[string]any{
		"aps": map[string]any{
			"alert": map[string]any{
				"title": title,
				"body":  body,
			},
			"sound": "default",
		},
	}

	// Adjuntar data custom fuera de "aps"
	for k, v := range data {
		apns[k] = v
	}

	pushMessage := dtos.PushMessage{
		Default: body,
		GCM:     gcm,
		APNS:    apns,
	}

	// publish options
	publishOptions := services.SnsPublishOptions{
		Subject: title,
	}

	if attributes != nil {
		attrs := services.SnsMessageAttributes{}

		for k, v := range attributes {
			attrs[k] = types.MessageAttributeValue{
				DataType:    aws.String("String"),
				StringValue: aws.String(v),
			}
		}

		publishOptions.Attributes = attrs
	}

	return pushMessage, publishOptions
}
