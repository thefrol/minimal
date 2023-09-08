package amazon

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/caarlos0/env/v9"
)

// Client создает клиента s3 из креденшалсов, хранящихся в переменных окружения
// MNML_KEY, MNMAL_SECRET, MNML_SESSION(опционально)
func Client() (*s3.Client, error) {
	e := struct {
		Key     string `env:"MNML_KEY"`
		Secret  string `env:"MNML_SECRET"`
		Session string `env:"MNML_SESSION"`
	}{}

	if err := env.Parse(&e); err != nil {
		fmt.Printf("%+v\n", err)
		return nil, err
	}

	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			PartitionID:   "yc",
			URL:           "https://storage.yandexcloud.net",
			SigningRegion: "ru-central1",
			//check setting another region
		}, nil
	})

	creds := func(cfg *config.LoadOptions) error {
		cfg.Credentials = credentials.NewStaticCredentialsProvider(e.Key, e.Secret, e.Session)
		return nil
	}

	// Подгружаем конфигрурацию
	cfg, err := config.LoadDefaultConfig(context.TODO(), creds, config.WithEndpointResolverWithOptions(customResolver))
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	// Создаем клиента для доступа к хранилищу S3
	client := s3.NewFromConfig(cfg)

	return client, nil

}

// код взяд из официальной документации yandex cloud
