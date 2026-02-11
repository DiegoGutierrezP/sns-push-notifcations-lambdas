package config

import (
	"log"
	"os"
	"sync"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type (
	Config struct {
		Sns               Sns      `envPrefix:"SNS_"`
		DynamoDb          DynamoDb `envPrefix:"DYNAMODB_"`
		Optimove          Optimove `envPrefix:"OPTIMOVE_"`
		ApiRateLimitTable string   `env:"API_RATE_LIMIT_TABLE"`
	}

	Sns struct {
		Region         string `env:"REGION,required"`
		PlatformAppArn string `env:"PLATFORM_APPLICATION_ARN,required"`
	}

	DynamoDb struct {
		Region string `env:"REGION,required"`

		DevicesTable              string `env:"DEVICES_TABLE,required"`
		SubscriptionTable         string `env:"SUBSCRIPTION_TABLE,required"`
		NotificationRequestsTable string `env:"NOTIFICATION_REQUESTS_TABLE,required"`
		NotificationLogsTable     string `env:"NOTIFICATION_LOGS_TABLE,required"`
	}

	Optimove struct {
		ApiUrl          string `env:"API_URL"`
		ApiKey          string `env:"API_KEY"`
		CallbackBaseUrl string `env:"CALLBACK_BASE_URL"`
	}
)

var (
	once           sync.Once
	configInstance Config
)

func GetConfig() *Config {
	once.Do(func() {
		//Intenta cargar .env SOLO si existe (en Lambda normalmente NO existe)
		if _, err := os.Stat(".env"); err == nil {
			if err := godotenv.Load(); err != nil {
				log.Printf("Warning: .env could not be loaded: %v", err)
			}
		} else {
			log.Println(".env not found. Continuing with environment variables.")
		}

		// Parsear variables de entorno reales (de Lambda)
		if err := env.Parse(&configInstance); err != nil {
			log.Panicf("Error parsing environment variables: %v", err)
		}
	})
	return &configInstance

}
