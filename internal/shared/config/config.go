package config

import (
	"sync"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type (
	Config struct {
		Db       Db       `envPrefix:"DB_"`
		Sns      Sns      `envPrefix:"SNS_"`
		DynamoDb DynamoDb `envPrefix:"DYNAMODB_"`
		Optimove Optimove `envPrefix:"OPTIMOVE_"`
	}

	Db struct {
		Host     string `env:"HOST"`
		Port     string `env:"PORT"`
		User     string `env:"USER"`
		Password string `env:"PASSWORD"`
		DBName   string `env:"DBNAME"`
		SSL      bool   `env:"SSL"`
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
	}

	Optimove struct {
		ApiUrl string `env:"API_URL"`
		ApiKey string `env:"API_KEY"`
	}
)

var (
	once           sync.Once
	configInstance Config
)

func GetConfig() *Config {
	once.Do(func() {
		if err := godotenv.Load(); err != nil {
			panic("Error cargando .env: " + err.Error())
		}

		if err := env.Parse(&configInstance); err != nil {
			panic(err)
		}
	})
	// fmt.Printf("%+v", configInstance)
	return &configInstance
}
