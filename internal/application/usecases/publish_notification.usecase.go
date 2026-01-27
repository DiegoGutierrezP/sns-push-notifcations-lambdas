package usecases

import (
	"context"
	"encoding/json"
	"errors"
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

	if rq.TargetArn == "" && rq.TopicArn == "" {
		return errors.New("targetArn or topicArn is required")
	}

	dataBytes, _ := json.Marshal(rq.Data)

	payload := map[string]string{
		"default": rq.Body,
		"GCM": fmt.Sprintf(`{
		"notification": {
			"title": "%s",
			"body": "%s"
		},
		"data": %s
	}`, rq.Title, rq.Body, string(dataBytes)),
	}

	attrs := services.SnsMessageAttributes{}

	for k, v := range rq.Attributes {
		attrs[k] = types.MessageAttributeValue{
			DataType:    aws.String("String"),
			StringValue: aws.String(v),
		}
	}

	publishOptions := services.SnsPublishOptions{
		Attributes: attrs,
		Subject:    rq.Title,
	}

	uc.snsService.Publish(ctx, &rq.TopicArn, &rq.TargetArn, payload, &publishOptions)

	return nil
}
