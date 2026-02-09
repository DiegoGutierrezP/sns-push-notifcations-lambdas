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
		Region          string `env:"REGION"`
		AccessKeyId     string `env:"ACCESS_KEY_ID"`
		SecretAccessKey string `env:"SECRET_ACCESS_KEY"`
		PlatformAppArn  string `env:"PLATFORM_APPLICATION_ARN"`
	}

	DynamoDb struct {
		Region          string `env:"REGION"`
		AccessKeyId     string `env:"ACCESS_KEY_ID"`
		SecretAccessKey string `env:"SECRET_ACCESS_KEY"`

		DevicesTable              string `env:"DEVICES_TABLE"`
		SubscriptionTable         string `env:"SUBSCRIPTION_TABLE"`
		NotificationRequestsTable string `env:"NOTIFICATION_REQUESTS_TABLE"`
		NotificationLogsTable     string `env:"NOTIFICATION_LOGS_TABLE"`
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
	// once.Do(func() {
	// 	if err := godotenv.Load(); err != nil {
	// 		panic("Error cargando .env: " + err.Error())
	// 	}

	// 	if err := env.Parse(&configInstance); err != nil {
	// 		panic(err)
	// 	}
	// })
	// // fmt.Printf("%+v", configInstance)
	// return &configInstance

	once.Do(func() {
		//Intenta cargar .env SOLO si existe (en Lambda normalmente NO existe)
		if _, err := os.Stat(".env"); err == nil {
			if err := godotenv.Load(); err != nil {
				log.Printf("Advertencia: no se pudo cargar .env: %v", err)
			}
		} else {
			log.Println("No se encontró .env. Continuando con variables de entorno.")
		}

		// Parsear variables de entorno reales (de Lambda)
		if err := env.Parse(&configInstance); err != nil {
			log.Panicf("Error parseando variables de entorno: %v", err)
		}
	})
	return &configInstance

}
