package database

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"

	appConfig "lmbd-digital-push-notifications/internal/shared/config"
)

func NewDynamoDbConnection(ctx context.Context, cfg *appConfig.Config) *dynamodb.Client {
	creds := aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(
		cfg.DynamoDb.AccessKeyId,
		cfg.DynamoDb.SecretAccessKey,
		"",
	))

	awsCfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(cfg.DynamoDb.Region),
		config.WithCredentialsProvider(creds),
	)

	if err != nil {
		log.Fatalf("error cargando configuración AWS: %v", err)
	}

	return dynamodb.NewFromConfig(awsCfg, func(o *dynamodb.Options) {})
}
