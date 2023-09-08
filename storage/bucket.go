package storage

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/thefrol/minimal/internal/amazon"
)

type Bucket struct {
	s3client *s3.Client
	Name     string
}

func New(name string) (*Bucket, error) {
	c, err := amazon.Client()
	if err != nil {
		return nil, err
	}
	b := Bucket{
		s3client: c,
		Name:     name,
	}
	return &b, nil
}

// UploadFIle загружает файл в бакет. objectkey - ключ объкта в бакете
func (b Bucket) UploadFile(fileName string, objectKey string) error {
	file, err := os.Open(fileName)
	if err != nil {
		log.Printf("Couldn't open file %v to upload. Here's why: %v\n", fileName, err)
	} else {
		defer file.Close()
		_, err = b.s3client.PutObject(context.TODO(), &s3.PutObjectInput{
			Bucket: aws.String(b.Name),
			Key:    aws.String(objectKey),
			Body:   file,
		})
		if err != nil {
			log.Printf("Couldn't upload file %v to %v:%v. Here's why: %v\n",
				fileName, b.Name, objectKey, err)
		}
	}
	return err
}

// Возвращает содержимое файла objectkey. передает поток, который требует закрытия
func (b Bucket) GetReader(objectKey string) (io.ReadCloser, error) {
	o, err := b.s3client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(b.Name),
		Key:    aws.String(objectKey),
	})
	if err != nil {
		fmt.Printf("Невозможно получить %v из бакета %v, по причине %+v", objectKey, b.Name, err)
		return nil, err
	}
	return o.Body, nil
}

// Возвращает содержимое файла objectkey, передает слайс байт
func (b Bucket) Get(objectKey string) ([]byte, error) {
	r, err := b.GetReader(objectKey)
	if err != nil {
		return nil, err
	}

	buf, err := io.ReadAll(r)
	if err != nil {
		fmt.Printf("Невозможно получить %v из бакета %v, по причине %+v", objectKey, b.Name, err)
		return nil, err
	}

	return buf, nil
}

// Возвращает содержимое файла objectkey, передает строку
func (b Bucket) GetString(objectKey string) (string, error) {
	buf, err := b.Get(objectKey)
	if err != nil {
		return "", err
	}
	return string(buf), err
}

// Возвращает список ключей бакета
func (b Bucket) Objects() ([]Object, error) {
	result, err := b.s3client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: aws.String(b.Name),
	})
	if err != nil {
		fmt.Printf("Невозможно получить список файлов %v, по причине %+v", b.Name, err)
		return nil, err
	}
	oo := []Object{}
	for _, o := range result.Contents {
		oo = append(oo, Object{
			Key:  *o.Key,
			Size: o.Size,
			Date: *o.LastModified,
		})
	}
	return oo, nil
}

// Names возвращает слайс имен объктов в бакете
func (b Bucket) Names() ([]string, error) {
	objects, err := b.Objects()
	if err != nil {
		return nil, err
	}
	sl := []string{}
	for _, o := range objects {
		sl = append(sl, o.Key)
	}
	return sl, nil
}

type Object struct {
	Key  string
	Size int64
	Date time.Time
}
