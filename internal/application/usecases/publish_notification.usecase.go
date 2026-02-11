package usecases

import (
	"context"
	"encoding/json"
	"lmbd-digital-push-notifications/internal/application/contracts/repositories"
	"lmbd-digital-push-notifications/internal/application/contracts/services"
	"lmbd-digital-push-notifications/internal/application/dtos"
	appErrors "lmbd-digital-push-notifications/internal/application/errors"
	"lmbd-digital-push-notifications/internal/domain/constants"
	domainErrors "lmbd-digital-push-notifications/internal/domain/errors"
	"log/slog"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns/types"
)

type IPublishNotificationUseCase interface {
	Execute(ctx context.Context, rq dtos.PublishNotificationRequest) (*dtos.PublishNotificationResponse, error)
}

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
) IPublishNotificationUseCase {
	return &PublishNotificationUseCase{
		snsService:             snsService,
		notificationRepository: notificationRepository,
		subscriptionRepository: subscriptionRepository,
		logger:                 logger,
	}
}

func (uc *PublishNotificationUseCase) Execute(ctx context.Context, rq dtos.PublishNotificationRequest) (*dtos.PublishNotificationResponse, error) {
	uc.logger.Info("PublishNotificationUseCase started:",
		"TopicArn", rq.TopicArn,
		"TargetArn", rq.TargetArn,
	)

	if rq.TargetArn == nil && rq.TopicArn == nil {
		uc.logger.Error("targetArn and topicArn are empty")

		// return nil, fmt.Errorf("targetArn or topicArn is required")
		return nil, appErrors.NewApplicationError(
			domainErrors.APIMissingParams,
			http.StatusBadRequest,
			"targetArn and topicArn are empty",
			nil,
		)
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

		// return nil, fmt.Errorf("An error occurred while publishing notification: %w", err)
		return nil, appErrors.NewApplicationError(
			domainErrors.SNSPublishFailed,
			http.StatusBadRequest,
			"An error occurred while publishing notification",
			err,
		)
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

	return &dtos.PublishNotificationResponse{
		MessageId:    *messageId,
		TotalDevices: &notificationRequestDto.TotalDevices,
	}, nil
}

func (uc *PublishNotificationUseCase) buildPushMessage(
	title *string,
	body string,
	data map[string]string,
	attributes map[string]string,
) (dtos.PushMessage, services.SnsPublishOptions) {

	// --- ANDROID (FCM v1) ---
	gcmMap := map[string]any{
		"fcmV1Message": map[string]any{
			"message": map[string]any{
				"notification": map[string]any{
					"title": title,
					"body":  body,
				},
				"data": data,
			},
		},
	}

	// --- iOS (APNS) ---
	apnsMap := map[string]any{
		"aps": map[string]any{
			"alert": map[string]any{
				"title": title,
				"body":  body,
			},
			"sound": "default",
		},
	}
	for k, v := range data {
		apnsMap[k] = v
	}

	// Serializar a STRING JSON (requisito SNS)
	gcmStr, _ := json.Marshal(gcmMap)
	apnsStr, _ := json.Marshal(apnsMap)

	pushMessage := dtos.PushMessage{
		Default: body,
		GCM:     string(gcmStr),  // <--- SNS exige STRING JSON
		APNS:    string(apnsStr), // <--- SNS exige STRING JSON
	}

	// Opciones de publicación
	opts := services.SnsPublishOptions{
		Subject: aws.String(*title),
		// Esto lo sobreescribe Publish() si targetArn != nil
		MessageStructure: aws.String("json"),
	}

	// Si hay atributos, los agregamos
	if attributes != nil {
		attrs := services.SnsMessageAttributes{}
		for k, v := range attributes {
			attrs[k] = types.MessageAttributeValue{
				DataType:    aws.String("String"),
				StringValue: aws.String(v),
			}
		}
		opts.Attributes = attrs
	}

	return pushMessage, opts
}
