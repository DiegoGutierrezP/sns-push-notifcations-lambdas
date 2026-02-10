package aws

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"lmbd-digital-push-notifications/internal/application/contracts/services"
	appConfig "lmbd-digital-push-notifications/internal/shared/config"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

type SnsService struct {
	client         *sns.Client
	platformAppArn string
}

func NewSnsService(cfg *appConfig.Config) (services.ISnsService, error) {
	creds := aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(
		cfg.Sns.AccessKeyId,
		cfg.Sns.SecretAccessKey,
		"",
	))

	awsCfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion(cfg.Sns.Region),
		config.WithCredentialsProvider(creds),
	)

	if err != nil {
		return nil, fmt.Errorf("error cargando configuración: %w", err)
	}

	client := sns.NewFromConfig(awsCfg)

	return &SnsService{
		client:         client,
		platformAppArn: cfg.Sns.PlatformAppArn,
	}, nil

}

func (s *SnsService) Subscription(
	ctx context.Context,
	topicArn string,
	endpoint string,
	protocol string,
	attributes *services.SnsSubscriptionAttributes,
) (*string, error) {
	cctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// var attributesMap map[string]string

	// if attributes != nil {
	// 	jsonData, err := json.Marshal(attributes)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 		return nil, err
	// 	}

	// 	err = json.Unmarshal(jsonData, &attributesMap)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// }

	// if protocol == "application" {
	// 	delete(attributesMap, "RawMessageDelivery")
	// }

	attrs, err := s.buildSnsSubscriptionAttributes(protocol, attributes)
	if err != nil {
		return nil, err
	}

	output, err := s.client.Subscribe(cctx, &sns.SubscribeInput{
		TopicArn:              aws.String(topicArn),
		Endpoint:              aws.String(endpoint),
		Protocol:              aws.String(protocol),
		Attributes:            attrs,
		ReturnSubscriptionArn: true,
	})

	if err != nil {
		return nil, err
	}

	return output.SubscriptionArn, nil
}

func (s *SnsService) Unsubscription(ctx context.Context, subscriptionArn string) error {
	cctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	_, err := s.client.Unsubscribe(cctx, &sns.UnsubscribeInput{
		SubscriptionArn: aws.String(subscriptionArn),
	})

	if err != nil {
		return err
	}

	return nil
}

func (s *SnsService) UpdateSubscriptionFilterPolicy(
	ctx context.Context,
	subscriptionArn string,
	newFilterPolicy map[string][]string,
) error {

	// 1. Validar FilterPolicy antes de enviarlo a AWS
	if err := s.validateFilterPolicy(newFilterPolicy); err != nil {
		return err
	}

	// 2. Serializar FilterPolicy a JSON
	jsonData, err := json.Marshal(newFilterPolicy)
	if err != nil {
		return fmt.Errorf("error serializando filter policy: %w", err)
	}

	filterPolicyStr := string(jsonData)

	// 3. Setear atributo en SNS
	cctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	_, err = s.client.SetSubscriptionAttributes(cctx, &sns.SetSubscriptionAttributesInput{
		SubscriptionArn: aws.String(subscriptionArn),
		AttributeName:   aws.String("FilterPolicy"),
		AttributeValue:  aws.String(filterPolicyStr),
	})

	if err != nil {
		return fmt.Errorf("error actualizando filter policy en SNS: %w", err)
	}

	return nil
}

func (s *SnsService) UpdateSubscriptionAttributes(
	ctx context.Context,
	subscriptionArn string,
	attributes *services.SnsSubscriptionAttributes,
) error {

	attrs, err := s.buildSnsSubscriptionAttributes("", attributes)
	if err != nil {
		return err
	}

	if attrs == nil {
		return fmt.Errorf("no se enviaron atributos para actualizar")
	}

	cctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	for key, val := range attrs {
		_, err := s.client.SetSubscriptionAttributes(cctx, &sns.SetSubscriptionAttributesInput{
			SubscriptionArn: aws.String(subscriptionArn),
			AttributeName:   aws.String(key),
			AttributeValue:  aws.String(val),
		})

		if err != nil {
			return fmt.Errorf("error actualizando atributo %s: %w", key, err)
		}
	}

	// v := reflect.ValueOf(attrs).Elem()
	// t := reflect.TypeOf(attrs).Elem()

	// for i := 0; i < t.NumField(); i++ {
	// 	field := t.Field(i)
	// 	value := v.Field(i)

	// 	tag := field.Tag.Get("sns")
	// 	if tag == "" {
	// 		continue // no es un atributo SNS
	// 	}

	// 	parts := strings.Split(tag, ",")
	// 	attrName := parts[0]
	// 	attrType := parts[1]

	// 	// si es puntero y nil → no se actualiza
	// 	if value.Kind() == reflect.Pointer && value.IsNil() {
	// 		continue
	// 	}

	// 	var attrValue string

	// 	switch attrType {
	// 	case "json":
	// 		jsonData, err := json.Marshal(value.Interface())
	// 		if err != nil {
	// 			return fmt.Errorf("error serializando %s: %w", attrName, err)
	// 		}

	// 		attrValue = string(jsonData)

	// 	case "bool":
	// 		b := value.Elem().Bool()
	// 		if b {
	// 			attrValue = "true"
	// 		} else {
	// 			attrValue = "false"
	// 		}

	// 	case "string":
	// 		attrValue = value.Elem().String()

	// 	default:
	// 		return fmt.Errorf("tipo no soportado: %s", attrType)
	// 	}

	// 	// enviar atributo a SNS
	// 	_, err := s.client.SetSubscriptionAttributes(cctx, &sns.SetSubscriptionAttributesInput{
	// 		SubscriptionArn: aws.String(subscriptionArn),
	// 		AttributeName:   aws.String(attrName),
	// 		AttributeValue:  aws.String(attrValue),
	// 	})

	// 	if err != nil {
	// 		return fmt.Errorf("error actualizando atributo %s: %w", attrName, err)
	// 	}
	// }

	return nil
}

func (s *SnsService) Publish(
	ctx context.Context,
	topicArn *string,
	targetArn *string,
	message any,
	opts *services.SnsPublishOptions,
) (*string, error) {

	if topicArn == nil && targetArn == nil {
		return nil, fmt.Errorf("debes enviar topicArn o targetArn")
	}

	if topicArn != nil && targetArn != nil {
		return nil, fmt.Errorf("solo uno de topicArn o targetArn puede estar presente")
	}

	input := &sns.PublishInput{}

	// Destino
	if topicArn != nil {
		input.TopicArn = topicArn
	}
	if targetArn != nil {
		input.TargetArn = targetArn
	}

	// Serializar mensaje
	switch m := message.(type) {
	case string:
		input.Message = aws.String(m)
		// Si es TOPIC, y viene como string, intentar activar MessageStructure solo si es un multi-protocolo válido
		if topicArn != nil {
			input.MessageStructure = aws.String("json")
		}
	default:
		b, err := json.Marshal(m)
		if err != nil {
			return nil, fmt.Errorf("error serializando json: %w", err)
		}
		input.Message = aws.String(string(b))

		// Para TOPIC siempre usamos mensaje multi‑protocolo
		if topicArn != nil {
			input.MessageStructure = aws.String("json")
		}
	}

	// Opciones adicionales
	if opts != nil {

		if opts.Subject != nil {
			input.Subject = aws.String(*opts.Subject)
		}

		if opts.Attributes != nil {
			input.MessageAttributes = opts.Attributes
		}

		if opts.MessageGroupId != "" {
			input.MessageGroupId = aws.String(opts.MessageGroupId)
		}

		if opts.MessageDedupId != "" {
			input.MessageDeduplicationId = aws.String(opts.MessageDedupId)
		}

		// Permitir override explícito
		if opts.MessageStructure != nil {
			input.MessageStructure = opts.MessageStructure
		}
	}

	// Publicar
	res, err := s.client.Publish(ctx, input)
	if err != nil {
		return nil, err
	}

	return res.MessageId, nil

}

func (s *SnsService) PublishToTopic(ctx context.Context, topicArn string, message any, opts *services.SnsPublishOptions) (*string, error) {
	return s.Publish(ctx, &topicArn, nil, message, opts)
}

func (s *SnsService) PublishToTarget(ctx context.Context, targetArn string, message any, opts *services.SnsPublishOptions) (*string, error) {
	return s.Publish(ctx, nil, &targetArn, message, opts)
}

func (s *SnsService) CreateEndpoint(ctx context.Context, deviceToken string) (string, error) {

	cctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	output, err := s.client.CreatePlatformEndpoint(cctx, &sns.CreatePlatformEndpointInput{
		PlatformApplicationArn: aws.String(s.platformAppArn),
		Token:                  aws.String(deviceToken),
		Attributes: map[string]string{
			"Enabled": "true",
		},
	})

	if err != nil {
		return "", fmt.Errorf("error creando endpoint SNS: %w", err)
	}

	return aws.ToString(output.EndpointArn), nil
}

func (s *SnsService) UpdateEndpoint(ctx context.Context, endpointArn string, newToken string) error {

	cctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	_, err := s.client.SetEndpointAttributes(cctx, &sns.SetEndpointAttributesInput{
		EndpointArn: aws.String(endpointArn),
		Attributes: map[string]string{
			"Token":   newToken,
			"Enabled": "true",
		},
	})

	return err
}

func (s *SnsService) buildSnsSubscriptionAttributes(protocol string, a *services.SnsSubscriptionAttributes) (map[string]string, error) {
	if a == nil {
		return nil, nil
	}

	attrs := map[string]string{}

	// RawMessageDelivery (solo para sqs/http/https/lambda)
	if a.RawMessageDelivery != nil && protocol != "application" {
		if *a.RawMessageDelivery {
			attrs["RawMessageDelivery"] = "true"
		} else {
			attrs["RawMessageDelivery"] = "false"
		}
	}

	// FilterPolicy (debe ser JSON string)
	if a.FilterPolicy != nil {
		b, err := json.Marshal(a.FilterPolicy)
		if err != nil {
			return nil, fmt.Errorf("filter policy json: %w", err)
		}
		attrs["FilterPolicy"] = string(b)
	}

	// FilterPolicyScope (solo valores válidos)
	if a.FilterPolicyScope != nil {
		v := *a.FilterPolicyScope
		if v == "MessageAttributes" || v == "MessageBody" {
			attrs["FilterPolicyScope"] = v
		}
	}

	// DeliveryPolicy (JSON string)
	if a.DeliveryPolicy != nil {
		b, err := json.Marshal(a.DeliveryPolicy)
		if err != nil {
			return nil, fmt.Errorf("delivery policy json: %w", err)
		}
		attrs["DeliveryPolicy"] = string(b)
	}

	// RedrivePolicy (JSON string)
	if a.RedrivePolicy != nil {
		b, err := json.Marshal(a.RedrivePolicy)
		if err != nil {
			return nil, fmt.Errorf("redrive policy json: %w", err)
		}
		attrs["RedrivePolicy"] = string(b)
	}

	if len(attrs) == 0 {
		return nil, nil
	}

	return attrs, nil
}

func (s *SnsService) validateFilterPolicy(filter map[string][]string) error {
	var maxFilterKeys int = 5
	var maxFilterCombinations int = 150
	var maxFilterPolicySize int = 256 * 1024

	// 1. Validar número máximo de keys
	if len(filter) > maxFilterKeys {
		return fmt.Errorf("filter policy excede el número máximo de %d keys (tiene %d)", maxFilterKeys, len(filter))
	}

	// 2. Validar combinaciones (producto cartesiano)
	totalCombinations := 1
	for key, values := range filter {
		// key no debe estar vacío
		if key == "" {
			return errors.New("filter policy contiene una key vacía")
		}
		// cada key debe tener al menos un valor
		if len(values) == 0 {
			return fmt.Errorf("la key '%s' no puede tener una lista vacía", key)
		}
		totalCombinations *= len(values)
	}

	if totalCombinations > maxFilterCombinations {
		return fmt.Errorf("el filter policy tiene %d combinaciones y excede el máximo permitido de %d",
			totalCombinations, maxFilterCombinations)
	}

	// 3. Validar tamaño del JSON
	jsonData, err := json.Marshal(filter)
	if err != nil {
		return fmt.Errorf("error serializando filter policy: %w", err)
	}

	if len(jsonData) > maxFilterPolicySize {
		return fmt.Errorf("filter policy supera el límite de tamaño (%d KB)", maxFilterPolicySize/1024)
	}

	return nil
}
