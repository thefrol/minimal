// Модуль для работы с бакетами. МОжно создавать, менять, удалять, получать файлы
//
// Бакет создается только нессколькими простыми формулами, из мерепенных окружения, если нужно что0то посложнее
// можно воспользоваться контруктами уровня модуля minimal, вроде
//	Profile("test").Bucket("dev")
//	DefaultProfile().Default().Bucket()
//	Тут у нас только из переменных окружения или из файлов

package bucket

import (
	"fmt"
	"strings"

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

	validate := not(emptyName, nameHasForbiddenSymbols, startsWithDigit)
	// по хорошему ошибки бы собирать и выводить их только если бакет не создастся
	if !validate(proto) {
		return nil, fmt.Errorf("неправильное имя бакета %v", proto.name)
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
	b, err := WithOptions(StaticCredentials(key, secret, session), WithName(name))
	return b, err
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

// validators

// На имя бакета накладываются следующие ограничения: #todo

// Длина имени должна быть от 3 до 63 символов.
// Имя может содержать строчные буквы латинского алфавита, цифры, дефисы и точки.
// Первый и последний символы должны быть буквами или цифрами.
// Символы справа и слева от точки должны быть буквами или цифрами.
// Имя не должно иметь вид IP-адреса (например 10.1.3.9).

type validator func(protoBucket) bool

func nameHasForbiddenSymbols(p protoBucket) bool {
	return strings.ContainsAny(p.name, "_&!@#$%^&*(),[]{}\"':;/\\")
}

func emptyName(p protoBucket) bool {
	return p.name == ""
}
func startsWithDigit(p protoBucket) bool {
	firstCharacter := p.name[0:1]
	r := strings.Contains("0123456789", firstCharacter)
	return r
}

func and(validators ...validator) validator {
	return func(pb protoBucket) bool {
		for _, v := range validators {
			if !v(pb) {
				return false
			}
		}
		return true
	}
}

func not(validators ...validator) validator {
	return func(pb protoBucket) bool {
		for _, v := range validators {
			if v(pb) {
				return false
			}
		}
		return true
	}
}

// #todo все эти функции выделить бы в отдельный файл с интерфейсом для AWS клиентировнных штук
