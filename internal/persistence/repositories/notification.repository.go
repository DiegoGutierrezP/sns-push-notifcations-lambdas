package repositories

import (
	"context"
	"fmt"
	"lmbd-digital-push-notifications/internal/application/contracts/repositories"
	"lmbd-digital-push-notifications/internal/persistence/models"
	"lmbd-digital-push-notifications/internal/shared/config"
	"lmbd-digital-push-notifications/internal/shared/utils"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
)

type NotificationRepository struct {
	tableRequest string
	tableLog     string
	client       *dynamodb.Client
	config       *config.Config
}

func NewNotificationRepository(config *config.Config, client *dynamodb.Client) repositories.INotificationRepository {
	return &NotificationRepository{
		tableRequest: config.DynamoDb.NotificationRequestsTable,
		tableLog:     config.DynamoDb.NotificationLogsTable,
		client:       client,
		config:       config,
	}
}

func (r *NotificationRepository) RegisterRequest(ctx context.Context, nr *repositories.NotificationRequestRegisterDto) (*string, error) {

	pk := uuid.New().String()
	now := time.Now().UTC()

	model := models.NotificationRequestModel{
		PK:           pk,
		MessageId:    nr.MessageId,
		TopicArn:     nr.TopicArn,
		TargetArn:    nr.TargetArn,
		Title:        nr.Title,
		Body:         nr.Body,
		Status:       fmt.Sprint(nr.Status),
		TotalDevices: nr.TotalDevices,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	item, err := attributevalue.MarshalMap(model)
	if err != nil {
		return nil, err
	}

	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName:           &r.tableRequest,
		Item:                item,
		ConditionExpression: aws.String("attribute_not_exists(pk)"),
	})

	return &pk, err
}

func (r *NotificationRepository) UpdateRequest(ctx context.Context, data *repositories.NotificationRequestUpdateDto) error {

	model := models.NotificationRequestModel{
		PK:           data.PK,
		MessageId:    data.MessageId,
		Status:       fmt.Sprint(data.Status),
		TotalDevices: *data.TotalDevices,
	}

	fields := utils.StructToMap(model, "dynamodbav")

	updateExpr := "SET "
	exprValues := map[string]types.AttributeValue{}
	exprNames := map[string]string{}

	i := 0
	for k, v := range fields {
		nameKey := "#f" + strconv.Itoa(i)
		valueKey := ":v" + strconv.Itoa(i)

		updateExpr += nameKey + " = " + valueKey + ","

		exprNames[nameKey] = k
		av, err := attributevalue.Marshal(v)
		if err != nil {
			return err
		}
		exprValues[valueKey] = av
		i++
	}

	updateExpr = strings.TrimSuffix(updateExpr, ",")

	_, err := r.client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: &r.tableRequest,
		Key: map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{Value: data.PK},
		},
		UpdateExpression:          aws.String(updateExpr),
		ExpressionAttributeNames:  exprNames,
		ExpressionAttributeValues: exprValues,
		ConditionExpression:       aws.String("attribute_exists(pk)"),
	})

	return err
}

func (r *NotificationRepository) RegisterLog(ctx context.Context, nl *repositories.NotificationLogRegisterDto) error {
	now := time.Now().UTC()

	model := models.NotificationLogModel{
		PK:             nl.MessageId,
		SK:             nl.EndpointArn,
		TopicArn:       nl.TopicArn,
		DeliveryId:     nl.DeliveryId,
		DeliveryStatus: nl.DeliveryStatus,
		StatusCode:     nl.StatusCode,
		Payload:        nl.Payload,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	item, err := attributevalue.MarshalMap(model)
	if err != nil {
		return err
	}

	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: &r.tableLog,
		Item:      item,
		// ConditionExpression: aws.String("attribute_not_exists(pk) AND attribute_not_exists(sk)"),
	})

	return err
}

func (r *NotificationRepository) UpsertLog(ctx context.Context, nl *repositories.NotificationLogRegisterDto) error {
	now := time.Now().UTC().Format(time.RFC3339Nano)

	_, err := r.client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: &r.tableLog,
		Key: map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{Value: nl.MessageId},
			"sk": &types.AttributeValueMemberS{Value: nl.EndpointArn},
		},

		// Upsert:
		// - createdAt se setea SOLO si no existe
		// - updatedAt siempre se actualiza
		UpdateExpression: aws.String(
			"SET topicArn = :topicArn, " +
				"deliveryId = :deliveryId, " +
				"deliveryStatus = :deliveryStatus, " +
				"statusCode = :statusCode, " +
				"payload = :payload, " +
				"createdAt = if_not_exists(createdAt, :createdAt), " +
				"updatedAt = :updatedAt",
		),

		ExpressionAttributeValues: map[string]types.AttributeValue{
			":topicArn":       &types.AttributeValueMemberS{Value: nl.TopicArn},
			":deliveryId":     &types.AttributeValueMemberS{Value: nl.DeliveryId},
			":deliveryStatus": &types.AttributeValueMemberS{Value: nl.DeliveryStatus},
			":statusCode":     &types.AttributeValueMemberN{Value: strconv.Itoa(nl.StatusCode)},
			":payload":        &types.AttributeValueMemberS{Value: nl.Payload},
			":createdAt":      &types.AttributeValueMemberS{Value: now},
			":updatedAt":      &types.AttributeValueMemberS{Value: now},
		},
	})

	return err
}
