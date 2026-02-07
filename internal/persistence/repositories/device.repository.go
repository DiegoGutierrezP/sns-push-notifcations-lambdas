package repositories

import (
	"context"
	"lmbd-digital-push-notifications/internal/application/contracts/repositories"
	"lmbd-digital-push-notifications/internal/domain/entities"
	"lmbd-digital-push-notifications/internal/persistence/mappers"
	"lmbd-digital-push-notifications/internal/persistence/models"
	"lmbd-digital-push-notifications/internal/shared/config"
	"lmbd-digital-push-notifications/internal/shared/utils"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type DeviceRepository struct {
	table  string
	client *dynamodb.Client
	config *config.Config
}

func NewDeviceRepository(config *config.Config, client *dynamodb.Client) repositories.IDeviceRepository {
	return &DeviceRepository{
		table:  config.DynamoDb.DevicesTable,
		client: client,
		config: config,
	}
}

func (r *DeviceRepository) Save(ctx context.Context, d *entities.DeviceEntity) error {

	model := mappers.ToDeviceModel(d)

	model.PlatformApplicationArn = r.config.Sns.PlatformAppArn

	item, err := attributevalue.MarshalMap(model)
	if err != nil {
		return err
	}

	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName:           &r.table,
		Item:                item,
		ConditionExpression: aws.String("attribute_not_exists(pk)"),
	})

	return err
}

func (r *DeviceRepository) GetByID(ctx context.Context, id string) (*entities.DeviceEntity, error) {

	out, err := r.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: &r.table,
		Key: map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{Value: id},
		},
	})

	if err != nil {
		return nil, err
	}

	if out.Item == nil {
		return nil, nil
	}

	var device models.DeviceModel
	if err := attributevalue.UnmarshalMap(out.Item, &device); err != nil {
		return nil, err
	}

	return mappers.ToDeviceEntity(&device), nil
}

func (r *DeviceRepository) GetByCalimacoId(ctx context.Context, calimacoId int) ([]entities.DeviceEntity, error) {

	out, err := r.client.Query(ctx, &dynamodb.QueryInput{
		TableName:              &r.table,
		IndexName:              aws.String(models.DeviceCalimacoIdIndex),
		KeyConditionExpression: aws.String("calimacoId = :v"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":v": &types.AttributeValueMemberS{Value: strconv.Itoa(calimacoId)},
		},
		//Limit: aws.Int32(1),
	})

	if err != nil {
		return nil, err
	}

	var devices []models.DeviceModel
	if err := attributevalue.UnmarshalListOfMaps(out.Items, &devices); err != nil {
		return nil, err
	}

	entities := utils.Map(devices, func(d models.DeviceModel) entities.DeviceEntity {
		return *mappers.ToDeviceEntity(&d)
	})

	return entities, nil
}

func (r *DeviceRepository) Delete(ctx context.Context, id string) error {

	_, err := r.client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: &r.table,
		Key: map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{Value: id},
		},
	})

	return err
}

func (r *DeviceRepository) UpdateStatus(ctx context.Context, id string, status int) error {
	_, err := r.client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: &r.table,
		Key: map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{Value: id},
		},
		UpdateExpression:         aws.String("SET #s = :status"),
		ExpressionAttributeNames: map[string]string{"#s": "status"},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":status": &types.AttributeValueMemberN{Value: strconv.Itoa(status)},
		},
	})

	return err
}

func (r *DeviceRepository) List(ctx context.Context) ([]*entities.DeviceEntity, error) {
	out, err := r.client.Scan(ctx, &dynamodb.ScanInput{
		TableName: &r.table,
	})

	if err != nil {
		return nil, err
	}

	var devices []*models.DeviceModel
	if err := attributevalue.UnmarshalListOfMaps(out.Items, &devices); err != nil {
		return nil, err
	}

	entities := utils.Map(devices, func(d *models.DeviceModel) *entities.DeviceEntity {
		return mappers.ToDeviceEntity(d)
	})

	return entities, nil
}

func (r *DeviceRepository) ExistsByToken(ctx context.Context, token string) (bool, error) {

	out, err := r.client.Query(ctx, &dynamodb.QueryInput{
		TableName:              &r.table,
		IndexName:              aws.String(models.DeviceDeviceTokenIndex),
		KeyConditionExpression: aws.String("deviceToken = :v"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":v": &types.AttributeValueMemberS{Value: token},
		},
		Limit: aws.Int32(1),
	})

	if err != nil {
		return false, err
	}

	return len(out.Items) > 0, nil
}

func (r *DeviceRepository) ExistsByTokenExcept(
	ctx context.Context,
	token string,
	exceptDeviceID string,
) (bool, error) {

	out, err := r.client.Query(ctx, &dynamodb.QueryInput{
		TableName: &r.table,
		IndexName: aws.String(models.DeviceDeviceTokenIndex),

		KeyConditionExpression: aws.String("deviceToken = :v"),

		FilterExpression: aws.String("pk <> :except"),

		ExpressionAttributeValues: map[string]types.AttributeValue{
			":v":      &types.AttributeValueMemberS{Value: token},
			":except": &types.AttributeValueMemberS{Value: exceptDeviceID},
		},

		Limit: aws.Int32(1),
	})

	if err != nil {
		return false, err
	}

	return len(out.Items) > 0, nil
}

func (r *DeviceRepository) Update(
	ctx context.Context,
	id string,
	data *entities.DeviceEntity,
) error {

	model := mappers.ToDeviceModel(data)

	fields := utils.StructToMap(model, "dynamodbav")

	delete(fields, "pk")
	delete(fields, "createdAt")

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
		TableName: &r.table,
		Key: map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{Value: id},
		},
		UpdateExpression:          aws.String(updateExpr),
		ExpressionAttributeNames:  exprNames,
		ExpressionAttributeValues: exprValues,
		ConditionExpression:       aws.String("attribute_exists(pk)"),
	})

	return err
}
