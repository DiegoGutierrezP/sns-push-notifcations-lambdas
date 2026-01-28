package repositories

import (
	"context"
	"errors"
	"lmbd-digital-push-notifications/internal/application/contracts/repositories"
	"lmbd-digital-push-notifications/internal/domain/entities"
	"lmbd-digital-push-notifications/internal/persistence/mappers"
	"lmbd-digital-push-notifications/internal/persistence/models"
	"lmbd-digital-push-notifications/internal/shared/utils"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type SubscriptionRepository struct {
	table  string
	client *dynamodb.Client
}

func NewSubscriptionRepository(client *dynamodb.Client) repositories.ISubscriptionRepository {
	return &SubscriptionRepository{
		table:  "pushnoti-subscriptions",
		client: client,
	}
}

func (r *SubscriptionRepository) Save(ctx context.Context, d *entities.SubscriptionEntity) error {
	model := mappers.ToSubscribeModel(d)

	// model.DeviceID = models.SubscriptionPk(model.DeviceID)
	// model.TopicID = models.SubscriptionSk(model.TopicID)

	item, err := attributevalue.MarshalMap(model)
	if err != nil {
		return err
	}

	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName:           &r.table,
		Item:                item,
		ConditionExpression: aws.String("attribute_not_exists(pk) AND attribute_not_exists(sk)"),
	})

	return err
}

func (r *SubscriptionRepository) Get(ctx context.Context, deviceId, topicId string) (*entities.SubscriptionEntity, error) {
	out, err := r.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: &r.table,
		Key: map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{Value: models.SubscriptionPk(deviceId)},
			"sk": &types.AttributeValueMemberS{Value: models.SubscriptionSk(topicId)},
		},
	})

	if err != nil {
		return nil, err
	}

	if out.Item == nil {
		return nil, errors.New("subscription not found")
	}

	var sub models.SubscriptionModel
	if err := attributevalue.UnmarshalMap(out.Item, &sub); err != nil {
		return nil, err
	}

	return mappers.ToSubscribeEntity(&sub), nil
}

func (r *SubscriptionRepository) Delete(ctx context.Context, deviceId, topicId string) error {
	_, err := r.client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: &r.table,
		Key: map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{Value: models.SubscriptionPk(deviceId)},
			"sk": &types.AttributeValueMemberS{Value: models.SubscriptionSk(topicId)},
		},
	})

	return err
}

func (r *SubscriptionRepository) ListByDevice(ctx context.Context, deviceId string) ([]*entities.SubscriptionEntity, error) {
	out, err := r.client.Query(ctx, &dynamodb.QueryInput{
		TableName:              &r.table,
		KeyConditionExpression: aws.String("pk = :pk"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{Value: models.SubscriptionPk(deviceId)},
		},
	})

	if err != nil {
		return nil, err
	}

	var subs []*models.SubscriptionModel
	if err := attributevalue.UnmarshalListOfMaps(out.Items, &subs); err != nil {
		return nil, err
	}

	entities := utils.Map(subs, func(d *models.SubscriptionModel) *entities.SubscriptionEntity {
		return mappers.ToSubscribeEntity(d)
	})

	return entities, nil
}

func (r *SubscriptionRepository) UpdateStatus(ctx context.Context, deviceId, topicId string, isActive bool) error {
	_, err := r.client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: &r.table,
		Key: map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{Value: models.SubscriptionPk(deviceId)},
			"sk": &types.AttributeValueMemberS{Value: models.SubscriptionSk(topicId)},
		},
		UpdateExpression: aws.String("SET isActive = :v"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":v": &types.AttributeValueMemberBOOL{Value: isActive},
		},
	})

	return err
}

func (r *SubscriptionRepository) UpdateAttributes(ctx context.Context, deviceId, topicId string, attributes string) error {
	_, err := r.client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: &r.table,
		Key: map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{Value: models.SubscriptionPk(deviceId)},
			"sk": &types.AttributeValueMemberS{Value: models.SubscriptionSk(topicId)},
		},
		UpdateExpression: aws.String("SET attributes = :v"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":v": &types.AttributeValueMemberS{Value: attributes},
		},
	})

	return err
}
