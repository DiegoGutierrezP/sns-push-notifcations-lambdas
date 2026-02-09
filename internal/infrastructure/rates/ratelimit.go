package rates

import (
	"context"
	"errors"
	"fmt"
	"lmbd-digital-push-notifications/internal/shared/config"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// RateLimiter provides a DynamoDB-backed fixed-window rate limiter.
// It tracks request counts per key and time window using an atomic UpdateItem
// with a conditional expression to enforce limits.
type RateLimiter struct {
	table  string
	client *dynamodb.Client
}

// NewRateLimiter creates a RateLimiter using the configured DynamoDB table.
func NewRateLimiter(client *dynamodb.Client, config *config.Config) *RateLimiter {
	return &RateLimiter{
		table:  config.ApiRateLimitTable,
		client: client,
	}
}

// Allow reports whether a request identified by key is allowed within the current window.
// It increments the counter atomically and returns false when the limit is reached.
// windowSeconds defines the fixed window size in seconds, and limit is the max count per window.
func (r *RateLimiter) Allow(ctx context.Context, key string, limit int, windowSeconds int64) (bool, error) {
	now := time.Now().Unix()
	window := now / windowSeconds

	pk := key
	sk := fmt.Sprintf("WINDOW#%d", window)

	ttl := (window+1)*windowSeconds + 5

	_, err := r.client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: aws.String(r.table),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: pk},
			"SK": &types.AttributeValueMemberS{Value: sk},
		},
		UpdateExpression:    aws.String("SET #c = if_not_exists(#c, :zero) + :one, expiresAt = :ttl"),
		ConditionExpression: aws.String("#c < :limit OR attribute_not_exists(#c)"),
		ExpressionAttributeNames: map[string]string{
			"#c": "count",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":zero":  &types.AttributeValueMemberN{Value: "0"},
			":one":   &types.AttributeValueMemberN{Value: "1"},
			":limit": &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", limit)},
			":ttl":   &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", ttl)},
		},
	})

	if err != nil {
		var cfe *types.ConditionalCheckFailedException
		if errors.As(err, &cfe) {
			return false, nil // rate limit hit
		}
		return false, err
	}

	return true, nil
}

// Wait blocks until a request for key is allowed within the current window.
// It retries across windows, sleeping until the next window boundary when the limit is hit.
// The call returns early if ctx is canceled.
func (r *RateLimiter) Wait(ctx context.Context, key string, limit int, windowSeconds int64) error {
	for {
		ok, err := r.Allow(ctx, key, limit, windowSeconds)
		if err != nil {
			return err
		}
		if ok {
			return nil
		}

		// Esperar hasta la próxima ventana
		now := time.Now().Unix()
		nextWindow := ((now / windowSeconds) + 1) * windowSeconds
		sleep := time.Duration(nextWindow-now) * time.Second

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(sleep):
		}
	}
}
