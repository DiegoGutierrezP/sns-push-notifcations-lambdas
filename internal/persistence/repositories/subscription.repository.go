package repositories

import (
	"context"
	"lmbd-digital-push-notifications/internal/application/contracts/repositories"
	"lmbd-digital-push-notifications/internal/domain/entities"
	"lmbd-digital-push-notifications/internal/persistence/mappers"
	"lmbd-digital-push-notifications/internal/persistence/models"
	"lmbd-digital-push-notifications/internal/shared/config"
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

func NewSubscriptionRepository(config *config.Config, client *dynamodb.Client) repositories.ISubscriptionRepository {
	return &SubscriptionRepository{
		table:  config.DynamoDb.SubscriptionTable,
		client: client,
	}
}

func (r *SubscriptionRepository) Save(ctx context.Context, d *entities.SubscriptionEntity) error {
	model := mappers.ToSubscribeModel(d)

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
		return nil, nil
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

func (r *SubscriptionRepository) DeleteByDevice(ctx context.Context, deviceId string) error {
	pk := models.SubscriptionPk(deviceId)

	// 1️⃣ Traer todos los SK del device
	out, err := r.client.Query(ctx, &dynamodb.QueryInput{
		TableName:              &r.table,
		KeyConditionExpression: aws.String("pk = :pk"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{Value: pk},
		},
		ProjectionExpression: aws.String("pk, sk"), // solo lo necesario
	})
	if err != nil {
		return err
	}

	// 2️⃣ No hay nada que borrar
	if len(out.Items) == 0 {
		return nil
	}

	// 3️⃣ Batch delete (máx 25 por request)
	const batchSize = 25
	for i := 0; i < len(out.Items); i += batchSize {
		end := i + batchSize
		if end > len(out.Items) {
			end = len(out.Items)
		}

		writeRequests := make([]types.WriteRequest, 0, batchSize)

		for _, item := range out.Items[i:end] {
			pkAttr := item["pk"]
			skAttr := item["sk"]

			writeRequests = append(writeRequests, types.WriteRequest{
				DeleteRequest: &types.DeleteRequest{
					Key: map[string]types.AttributeValue{
						"pk": pkAttr,
						"sk": skAttr,
					},
				},
			})
		}

		_, err := r.client.BatchWriteItem(ctx, &dynamodb.BatchWriteItemInput{
			RequestItems: map[string][]types.WriteRequest{
				r.table: writeRequests,
			},
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *SubscriptionRepository) CountByTopic(ctx context.Context, topicID string) (int32, error) {

	out, err := r.client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(r.table),
		IndexName:              aws.String(models.SubscriptionTopicIndex),
		KeyConditionExpression: aws.String("#topic = :topicId"),
		ExpressionAttributeNames: map[string]string{
			"#topic": "sk", // TopicID vive en "sk"
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":topicId": &types.AttributeValueMemberS{Value: models.SubscriptionSk(topicID)},
		},
		Select: types.SelectCount, // solo Count, no Items
	})
	if err != nil {
		return 0, err
	}
	return out.Count, nil
}
