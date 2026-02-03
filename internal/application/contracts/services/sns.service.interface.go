package services

import (
	"context"

	snsTypes "github.com/aws/aws-sdk-go-v2/service/sns/types"
)

type SnsSubscriptionAttributes struct {
	FilterPolicy       map[string][]string `sns:"FilterPolicy,json"` // MAX 150
	RawMessageDelivery *bool               `sns:"RawMessageDelivery,bool"`
	DeliveryPolicy     map[string]any      `sns:"DeliveryPolicy,json"`
	RedrivePolicy      map[string]any      `sns:"RedrivePolicy,json"`
	FilterPolicyScope  *string             `sns:"FilterPolicyScope,string"`
}

type SnsMessageAttributes map[string]snsTypes.MessageAttributeValue

type SnsPublishOptions struct {
	Subject        *string
	Attributes     SnsMessageAttributes
	MessageGroupId string // para FIFO
	MessageDedupId string // para FIFO
	Structure      string // json | string
}

type ISnsService interface {
	Subscription(ctx context.Context, topicArn string, endpoint string, protocol string, attributes *SnsSubscriptionAttributes) (*string, error)
	Unsubscription(ctx context.Context, subscriptionArn string) error
	UpdateSubscriptionAttributes(ctx context.Context, subscriptionArn string, attrs *SnsSubscriptionAttributes) error
	UpdateSubscriptionFilterPolicy(ctx context.Context, subscriptionArn string, newFilterPolicy map[string][]string) error
	Publish(ctx context.Context, topicArn *string, targetArn *string, message any, opts *SnsPublishOptions) error
	PublishToTopic(ctx context.Context, topicArn string, message any, opts *SnsPublishOptions) error
	PublishToTarget(ctx context.Context, targetArn string, message any, opts *SnsPublishOptions) error
	CreateEndpoint(ctx context.Context, deviceToken string) (string, error)
	UpdateEndpoint(ctx context.Context, endpointArn string, newToken string) error
}
