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
func Client(functOpts ...func(opts *config.LoadOptions) error) (*s3.Client, error) {

	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			PartitionID:   "yc",
			URL:           "https://storage.yandexcloud.net",
			SigningRegion: "ru-central1",
			//check setting another region
		}, nil
	})

	functOpts = append(functOpts, config.WithEndpointResolverWithOptions(customResolver))

	// Подгружаем конфигрурацию
	cfg, err := config.LoadDefaultConfig(context.TODO(), functOpts...)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	// Создаем клиента для доступа к хранилищу S3
	client := s3.NewFromConfig(cfg)

	return client, nil

}

// Загружает креденшался из переменных окружения MNML_KEY, MNML_SECRET
func LoadFromEnv(opts *config.LoadOptions) error {
	e := struct {
		Key     string `env:"MNML_KEY"`
		Secret  string `env:"MNML_SECRET"`
		Session string `env:"MNML_SESSION"`
	}{}

	if err := env.Parse(&e); err != nil {
		fmt.Printf("%+v\n", err)
		return err
	}

	opts.Credentials = credentials.NewStaticCredentialsProvider(e.Key, e.Secret, e.Session)
	return nil
}

func LoadFromFile(opts *config.LoadOptions) error {
	/* 	e := struct {
		Key     string `env:"MNML_KEY"`
		Secret  string `env:"MNML_SECRET"`
		Session string `env:"MNML_SESSION"`
	}{} */

	return fmt.Errorf("Client.LoadFromFile() не воплощена")
}

func StaticKeys(key, secret, session string) config.LoadOptionsFunc {
	return func(opts *config.LoadOptions) error {
		opts.Credentials = credentials.NewStaticCredentialsProvider(key, secret, session)
		return nil
	}
}

// код взяд из официальной документации yandex cloud
