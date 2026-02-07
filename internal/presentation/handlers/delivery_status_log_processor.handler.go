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
	"log"
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

type DeliveryStatusLogProcessorHandler struct {
	usecase *usecases.SaveDeliveryStatusLogUseCase
}

func NewDeliveryStatusLogProcessorHandler(
	usecase *usecases.SaveDeliveryStatusLogUseCase,
) *DeliveryStatusLogProcessorHandler {
	return &DeliveryStatusLogProcessorHandler{
		usecase: usecase,
	}
}

func (h *DeliveryStatusLogProcessorHandler) Handler(
	ctx context.Context,
	evt events.CloudwatchLogsEvent,
) error {

	// 1) decode base64 + gzip
	raw, err := base64.StdEncoding.DecodeString(evt.AWSLogs.Data)
	if err != nil {
		return fmt.Errorf("base64 decode awslogs.data: %w", err)
	}

	unzipped, err := gunzip(raw)
	if err != nil {
		return fmt.Errorf("gunzip awslogs.data: %w", err)
	}

	// 2) unmarshal envelope CloudWatch Logs
	var payload cwLogsPayload
	if err := json.Unmarshal(unzipped, &payload); err != nil {
		return fmt.Errorf("unmarshal cloudwatch logs payload: %w", err)
	}

	// 3) process each logEvent.message (SNS delivery JSON)
	for _, le := range payload.LogEvents {

		fmt.Printf("Processing logEventId=%s message=%s\n", le.ID, le.Message)

		var dl dtos.SnsDeliveryLog
		if err := json.Unmarshal([]byte(le.Message), &dl); err != nil {
			log.Printf("skip: invalid sns delivery json, logEventId=%s err=%v message=%q", le.ID, err, le.Message)
			continue
		}

		occurredAt := time.UnixMilli(le.Timestamp)

		in := dtos.SaveDeliveryStatusInput{
			MessageID:      dl.Notification.MessageID,
			EndpointArn:    dl.Delivery.Destination,
			TopicArn:       dl.Notification.TopicArn,
			DeliveryStatus: dl.Status,
			StatusCode:     dl.Delivery.StatusCode,
			Payload:        le.Message, // guarda el JSON raw si quieres
			OccurredAt:     occurredAt,
		}

		if err := h.usecase.Execute(ctx, in); err != nil {

			log.Printf("failed to persist logEventId=%s messageId=%s endpoint=%s err=%v",
				le.ID, in.MessageID, in.EndpointArn, err)

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
