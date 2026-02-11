package database

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"

	appConfig "lmbd-digital-push-notifications/internal/shared/config"
)

func NewDynamoDbConnection(cfg *appConfig.Config) *dynamodb.Client {
	awsCfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion(cfg.DynamoDb.Region),
	)

	if err != nil {
		log.Fatalf("error cargando configuración AWS: %v", err)
	}

	return dynamodb.NewFromConfig(awsCfg, func(o *dynamodb.Options) {})
}
