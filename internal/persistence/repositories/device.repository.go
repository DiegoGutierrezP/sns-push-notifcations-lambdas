package repositories

import (
	"context"
	"errors"
	"lmbd-digital-push-notifications/internal/application/contracts/repositories"
	"lmbd-digital-push-notifications/internal/domain/entities"
	"lmbd-digital-push-notifications/internal/persistence/mappers"
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
)

type DeviceRepository struct {
	table  string
	client *dynamodb.Client
	config *config.Config
}

func NewDeviceRepository(config *config.Config, client *dynamodb.Client) repositories.IDeviceRepository {
	return &DeviceRepository{
		table:  "pushnoti-devices",
		client: client,
		config: config,
	}
}

func (r *DeviceRepository) Save(ctx context.Context, d *entities.DeviceEntity) error {

	model := mappers.ToDeviceModel(d)

	model.GSI1PK = models.DeviceGSI1Pk(model.DeviceToken)
	model.GSI2PK = models.DeviceGSI2Pk(model.CalimacoId)
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

	pk := models.DevicePk(id)

	out, err := r.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: &r.table,
		Key: map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{Value: pk},
		},
	})

	if err != nil {
		return nil, err
	}

	if out.Item == nil {
		return nil, errors.New("device not found")
	}

	var device models.DeviceModel
	if err := attributevalue.UnmarshalMap(out.Item, &device); err != nil {
		return nil, err
	}

	return mappers.ToDeviceEntity(&device), nil
}

func (r *DeviceRepository) GetByCalimacoId(ctx context.Context, calimacoId string) ([]entities.DeviceEntity, error) {

	gsi2pk := models.DeviceGSI2Pk(calimacoId)

	out, err := r.client.Query(ctx, &dynamodb.QueryInput{
		TableName:              &r.table,
		IndexName:              aws.String("CalimacoIdIndex"),
		KeyConditionExpression: aws.String("gsi2pk = :v"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":v": &types.AttributeValueMemberS{Value: gsi2pk},
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

	pk := models.DevicePk(id)

	_, err := r.client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: &r.table,
		Key: map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{Value: pk},
		},
	})

	return err
}

func (r *DeviceRepository) UpdateStatus(ctx context.Context, id string, status int) error {
	_, err := r.client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: &r.table,
		Key: map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{Value: models.DevicePk(id)},
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

	gsi1pk := models.DeviceGSI1Pk(token)

	out, err := r.client.Query(ctx, &dynamodb.QueryInput{
		TableName:              &r.table,
		IndexName:              aws.String("TokenIndex"),
		KeyConditionExpression: aws.String("gsi1pk = :v"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":v": &types.AttributeValueMemberS{Value: gsi1pk},
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

	gsi1pk := models.DeviceGSI1Pk(token)
	exceptPK := models.DevicePk(exceptDeviceID)

	out, err := r.client.Query(ctx, &dynamodb.QueryInput{
		TableName: &r.table,
		IndexName: aws.String("TokenIndex"),

		KeyConditionExpression: aws.String("gsi1pk = :v"),

		FilterExpression: aws.String("pk <> :except"),

		ExpressionAttributeValues: map[string]types.AttributeValue{
			":v":      &types.AttributeValueMemberS{Value: gsi1pk},
			":except": &types.AttributeValueMemberS{Value: exceptPK},
		},

		Limit: aws.Int32(1),
	})

	if err != nil {
		return false, err
	}

	return len(out.Items) > 0, nil
}

func (r *DeviceRepository) UpdateFields(
	ctx context.Context,
	id string,
	fields map[string]any,
) error {

	if len(fields) == 0 {
		return nil
	}

	pk := models.DevicePk(id)

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
			"pk": &types.AttributeValueMemberS{Value: pk},
		},
		UpdateExpression:          aws.String(updateExpr),
		ExpressionAttributeNames:  exprNames,
		ExpressionAttributeValues: exprValues,
		ConditionExpression:       aws.String("attribute_exists(pk)"),
	})

	return err
}

func (r *DeviceRepository) Update(
	ctx context.Context,
	id string,
	update *repositories.UpdateDeviceFields,
) error {

	fields := utils.StructToMap(update)

	if len(fields) == 0 {
		return nil
	}

	fields["updatedAt"] = time.Now().UTC()

	delete(fields, "pk")
	delete(fields, "deviceToken")
	delete(fields, "gsi1pk")

	pk := models.DevicePk(id)

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
			"pk": &types.AttributeValueMemberS{Value: pk},
		},
		UpdateExpression:          aws.String(updateExpr),
		ExpressionAttributeNames:  exprNames,
		ExpressionAttributeValues: exprValues,
		ConditionExpression:       aws.String("attribute_exists(pk)"),
	})

	return err
}
