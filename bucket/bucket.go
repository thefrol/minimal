package bucket

import (
	"context"
	"errors"
	"io"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

// Bucket позволяет проводит операции над файлами в бакете
type Bucket struct {
	s3client *s3.Client
	Name     string
}

// UploadFIle загружает файл в бакет. objectkey - ключ объкта в бакете
func (b Bucket) Put(r io.Reader, objectKey string) error {
	_, err := b.s3client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(b.Name),
		Key:    aws.String(objectKey),
		Body:   r,
	})

	return err
}

// UploadFIle загружает файл в бакет. objectkey - ключ объкта в бакете
func (b Bucket) UploadFile(fileName string, objectKey string) error {
	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	err = b.Put(file, objectKey)
	return err
}

// Возвращает содержимое файла objectkey. передает поток, который требует закрытия
func (b Bucket) Get(objectKey string) (io.ReadCloser, error) {
	o, err := b.s3client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(b.Name),
		Key:    aws.String(objectKey),
	})
	if err != nil {
		var nsk *types.NoSuchKey
		if errors.As(err, &nsk) {
			return nil, &KeyNotFound{BucketName: b.Name, Key: objectKey, Err: err}
		}
		return nil, err
	}
	return o.Body, err
}

// Возвращает содержимое файла objectkey, передает слайс байт
func (b Bucket) GetBytes(objectKey string) ([]byte, error) {
	r, err := b.Get(objectKey)
	if err != nil {
		return nil, err
	}

	buf, err := io.ReadAll(r)
	return buf, err
}

// Возвращает содержимое файла objectkey, передает строку
func (b Bucket) GetString(objectKey string) (string, error) {
	r, err := b.Get(objectKey)
	if err != nil {
		return "", err
	}

	buf := new(strings.Builder)
	_, err = io.Copy(buf, r)
	return buf.String(), err
}

// Возвращает список ключей бакета
func (b Bucket) Objects() ([]Object, error) {
	result, err := b.s3client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: aws.String(b.Name),
	})
	if err != nil {
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

// Delete удаляет ключ objectName из бакета
func (b Bucket) Delete(objectKey string) error {
	_, err := b.s3client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(b.Name),
		Key:    aws.String(objectKey),
	})
	if err != nil {
		var nsk *types.NoSuchKey
		if errors.As(err, &nsk) {
			return &KeyNotFound{BucketName: b.Name, Key: objectKey, Err: err}
		}
		return err
	}
	return nil
}

type Object struct {
	Key  string
	Size int64
	Date time.Time
}
