package handlers

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"lmbd-digital-push-notifications/internal/application/dtos"
	"lmbd-digital-push-notifications/internal/application/usecases"
	"log/slog"
	"time"

	"github.com/aws/aws-lambda-go/events"
)

//Dtos

type cwLogsPayload struct {
	MessageType         string           `json:"messageType"`
	Owner               string           `json:"owner"`
	LogGroup            string           `json:"logGroup"`
	LogStream           string           `json:"logStream"`
	SubscriptionFilters []string         `json:"subscriptionFilters"`
	LogEvents           []cwLogsLogEvent `json:"logEvents"`
}

type cwLogsLogEvent struct {
	ID        string `json:"id"`
	Timestamp int64  `json:"timestamp"`
	Message   string `json:"message"`
}

// Handler

type DeliveryLogProcessorHandler struct {
	usecase usecases.ISaveDeliveryStatusLogUseCase
	logger  *slog.Logger
}

func NewDeliveryLogProcessorHandler(
	usecase usecases.ISaveDeliveryStatusLogUseCase,
	logger *slog.Logger,
) *DeliveryLogProcessorHandler {
	return &DeliveryLogProcessorHandler{
		usecase: usecase,
		logger:  logger,
	}
}

func (h *DeliveryLogProcessorHandler) Handler(
	ctx context.Context,
	evt events.CloudwatchLogsEvent,
) error {

	// 1) decode base64 + gzip
	raw, err := base64.StdEncoding.DecodeString(evt.AWSLogs.Data)
	if err != nil {
		h.logger.Error("base64 decode awslogs.data failed",
			"err", err,
		)
		return fmt.Errorf("base64 decode awslogs.data: %w", err)
	}

	unzipped, err := gunzip(raw)
	if err != nil {
		h.logger.Error("gunzip awslogs.data failed",
			"err", err,
		)
		return fmt.Errorf("gunzip awslogs.data: %w", err)
	}

	// 2) unmarshal envelope CloudWatch Logs
	var payload cwLogsPayload
	if err := json.Unmarshal(unzipped, &payload); err != nil {
		h.logger.Error("unmarshal cloudwatch logs payload failed",
			"err", err,
		)
		return fmt.Errorf("unmarshal cloudwatch logs payload: %w", err)
	}

	// 3) process each logEvent.message (SNS delivery JSON)
	for _, le := range payload.LogEvents {
		h.logger.Info("Processiong log event",
			"logEventId", le.ID,
			"message", le.Message,
		)

		var dl dtos.SnsDeliveryLog
		if err := json.Unmarshal([]byte(le.Message), &dl); err != nil {
			h.logger.Error("skip: invalid sns delivery json",
				"logEventId", le.ID,
				"message", le.Message,
				"err", err,
			)
			continue
		}

		occurredAt := time.UnixMilli(le.Timestamp)

		in := dtos.SaveDeliveryStatusInput{
			MessageID:      dl.Notification.MessageID,
			EndpointArn:    dl.Delivery.Destination,
			TopicArn:       dl.Notification.TopicArn,
			DeliveryStatus: dl.Status,
			StatusCode:     dl.Delivery.StatusCode,
			Payload:        le.Message,
			OccurredAt:     occurredAt,
		}

		if err := h.usecase.Execute(ctx, in); err != nil {

			h.logger.Error("failed to persist log",
				"logEventId", le.ID,
				"message", le.Message,
				"enpoint", in.EndpointArn,
				"err", err,
			)

			continue
		}

	}

	return nil

}

func gunzip(b []byte) ([]byte, error) {
	r, err := gzip.NewReader(bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	defer r.Close()

	out, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return out, nil
}
