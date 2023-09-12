// Модуль для работы с бакетами. МОжно создавать, менять, удалять, получать файлы
//
// Бакет создается только нессколькими простыми формулами, из мерепенных окружения, если нужно что0то посложнее
// можно воспользоваться контруктами уровня модуля minimal, вроде
//	Profile("test").Bucket("dev")
//	DefaultProfile().Default().Bucket()
//	Тут у нас только из переменных окружения или из файлов

package bucket

import (
	"errors"
	"fmt"

	"github.com/caarlos0/env"
	"github.com/thefrol/minimal/internal/amazon"
)

// New Создает новый объект бакета, и дает доступ в бакет с именем name
// Ключ и секрет берутся из переменных окружения по умолчанию, все
// остальное настраивается при помощи func_opts
//
// BUCKET_KEY, BUCKET_SECRET - переменные окружения с креденшансами
// BUCKET_BUCKET - имя бакета
func WithOptions(funcOpts ...configFunc) (*Bucket, error) {
	proto := protoBucket{}
	for _, f := range funcOpts {
		err := f(&proto)
		if err != nil {
			return nil, err
		}
	}

	// по хорошему ошибки бы собирать и выводить их только если бакет не создастся
	if proto.name == "" { //validate name! #todo без подчеркиваний там
		return nil, fmt.Errorf("пустое имя бакета")
	}

	c, err := amazon.Client(amazon.StaticKeys(proto.key, proto.secret, ""))
	if err != nil {
		return nil, err
	}

	b := new(Bucket)
	b.Name = proto.name
	b.s3client = c
	return b, nil
}

// New Создает новый объект бакета, и дает доступ в бакет с именем name
// Ключ и секрет берутся из переменных окружения по умолчанию, все
// остальное настраивается при помощи func_opts
//
// BUCKET_KEY, BUCKET_SECRET - переменные окружения с креденшансами
// BUCKET_BUCKET - имя бакета
func New(name string) (b *Bucket, err error) {
	b, err = WithOptions( /* ConfigFromFile(defaultConfigFile), */ CredentialsFromEnv, WithName(name))
	return
}

// Default открывает бакет, с креденшалсами по умолчанию(из переменных
// окружения), или из файла профиля, имя берется тоже оттуда же
func Default() (*Bucket, error) {
	b, err := WithOptions(CredentialsFromEnv) // тут может быть другой
	return b, err
}

// FromEnvironmentVariables Открывает бакет, полностью обусловленный переменными окружения,
// в том числе и имя бакета тоже берется из BUCKET_BUCKET,
// BUCKET_KEY, BUCKET_SECRET - тут креденшалсы
func FromEnvironmentVariables() (*Bucket, error) {
	b, err := WithOptions(CredentialsFromEnv) // #todo а еще добавить в цепочку
	return b, err
}

// FromKeys Открывает бакет, так что креденшалы задаются в открытом виде
func FromKeys(key, secret, session, name string) (*Bucket, error) {
	return nil, errors.New("функция FromFile Variables пока не воплощена")
}

// Опцефункции

type OptionsFunc func(b *Bucket) error
type configFunc func(b *protoBucket) error

func WithName(name string) configFunc {
	return func(proto *protoBucket) error {
		if name == "" {
			return nil // если имя пусток ничего не присваиваем дополнительно, может быть где-то в цепочке уже кто-то что-то ввел в имя
		}
		proto.name = name
		return nil
	}
}

func CredentialsFromEnv(proto *protoBucket) error {
	e := struct {
		Name   string `env:"BUCKET_NAME"`
		Key    string `env:"BUCKET_KEY"`
		Secret string `env:"BUCKET_SECRET"`
	}{}

	if err := env.Parse(&e); err != nil {
		fmt.Printf("%+v\n", err)
		return err
	}

	proto.name = e.Name
	proto.key = e.Key
	proto.secret = e.Secret
	return nil
}

func StaticCredentials(key, secret, name string) configFunc {
	return func(proto *protoBucket) error {
		proto.key = key
		proto.secret = secret
		proto.name = name
		return nil
	}
}

type protoBucket struct {
	name   string
	key    string
	secret string
}

// #todo все эти функции выделить бы в отдельный файл с интерфейсом для AWS клиентировнных штук
