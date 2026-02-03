package usecases

import (
	"context"
	"fmt"
	"lmbd-digital-push-notifications/internal/application/contracts/services"
	"lmbd-digital-push-notifications/internal/application/dtos"
	"log/slog"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns/types"
)

type PublishNotificationUseCase struct {
	snsService services.ISnsService
	logger     *slog.Logger
}

func NewPublishNotificationUseCase(snsService services.ISnsService, logger *slog.Logger) *PublishNotificationUseCase {
	return &PublishNotificationUseCase{
		snsService: snsService,
		logger:     logger,
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

	pushMessage := uc.buildPushMessage(rq.Title, rq.Body, rq.Data)

	publishOptions := services.SnsPublishOptions{
		Subject: rq.Title,
	}

	if rq.Attributes != nil {
		attrs := services.SnsMessageAttributes{}

		for k, v := range rq.Attributes {
			attrs[k] = types.MessageAttributeValue{
				DataType:    aws.String("String"),
				StringValue: aws.String(v),
			}
		}

		publishOptions.Attributes = attrs
	}

	if err := uc.snsService.Publish(ctx, rq.TopicArn, rq.TargetArn, pushMessage, &publishOptions); err != nil {
		uc.logger.Error("snsService.Publish failed",
			"TopicArn", rq.TopicArn,
			"TargetArn", rq.TargetArn,
			"err", err,
		)
		return fmt.Errorf("An error occurred while publishing notification: %w", err)
	}

	uc.logger.Info("Notification published successfully",
		"TopicArn", rq.TopicArn,
		"TargetArn", rq.TargetArn,
	)

	return nil
}

func (uc *PublishNotificationUseCase) buildPushMessage(
	title *string,
	body string,
	data map[string]string,
) dtos.PushMessage {

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

	return dtos.PushMessage{
		Default: body,
		GCM:     gcm,
		APNS:    apns,
	}
}
