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
	"github.com/thefrol/minimal/internal/amazon"
)

// Bucket позволяет проводит операции над файлами в бакете
type Bucket struct {
	s3client *s3.Client
	Name     string
}

// New Создает новый объект бакета, и дает доступ в бакет с именем name
// Ключ и секрет берутся из переменных окружения по умолчанию, все
// остальное настраивается при помощи func_opts
func New() (*Bucket, error) {
	c, err := amazon.Client()
	if err != nil {
		return nil, err
	}
	b := Bucket{
		s3client: c,
	}
	return &b, nil
}

// WithName открывает бакет, с креденшалсами по умолчанию(из переменных
// окружения), но можно настроить имя
func WithName(name string) (*Bucket, error) {
	b, err := New()
	if err != nil {
		return nil, err
	}
	b.Name = name //#todo переделать с func_opts
	return b, err
}

// FromEnvironmentVariables Открывает бакет, полностью обусловленный переменными окружения,
// в том числе и имя бакета тоже берется из
func FromEnvironmentVariables() (*Bucket, error) {
	return nil, errors.New("функция FromEnvironment Variables пока не воплощена")
}

// FromFile Открывает бакет, описанный в адресе,
// в том числе и имя бакета тоже берется из
func FromFile(f string) (*Bucket, error) {
	return nil, errors.New("функция FromFile Variables пока не воплощена")
}

// FromDefaultFile Открывает бакет, описанный в файле по умолчанию .minimal,
// в том числе и имя бакета тоже берется из
func FromDefaultFile() (*Bucket, error) {
	return nil, errors.New("функция FromFile Variables пока не воплощена")
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

type Object struct {
	Key  string
	Size int64
	Date time.Time
}
