package bucket

import (
	"errors"

	"github.com/thefrol/minimal/internal/amazon"
)

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
