package bucket

import (
	"errors"
	"fmt"
	"os"

	"github.com/thefrol/minimal/internal/amazon"
)

const (
	defaultConfigFile = ".minimal"
)

// New Создает новый объект бакета, и дает доступ в бакет с именем name
// Ключ и секрет берутся из переменных окружения по умолчанию, все
// остальное настраивается при помощи func_opts
//
// MNML_KEY, MNML_SECRET - переменные окружения с креденшансами
// MNML_BUCKET - имя бакета
func WithOptions(funcOpts ...OptionsFunc) (*Bucket, error) {
	b := new(Bucket)
	for _, f := range funcOpts {
		err := f(b)
		if err != nil {
			return nil, err
		}
	}
	// по хорошему ошибки бы собирать и выводить их только если бакет не создастся
	if b.Name == "" { //validate name! #todo без подчеркиваний там
		return nil, fmt.Errorf("пустое имя бакета")
	}
	return b, nil
}

// New Создает новый объект бакета, и дает доступ в бакет с именем name
// Ключ и секрет берутся из переменных окружения по умолчанию, все
// остальное настраивается при помощи func_opts
//
// MNML_KEY, MNML_SECRET - переменные окружения с креденшансами
// MNML_BUCKET - имя бакета
func New(name string) (b *Bucket, err error) {
	b, err = WithOptions( /* ConfigFromFile(defaultConfigFile), */ CredentialsFromEnv, NameFromEnv, WithName(name))
	return
}

// Default открывает бакет, с креденшалсами по умолчанию(из переменных
// окружения), или из файла профиля, имя берется тоже оттуда же
func Default() (*Bucket, error) {
	b, err := WithOptions(ConfigFromFile(defaultConfigFile), CredentialsFromEnv, NameFromEnv) // #todo а еще добавить в цепочку блять из файла
	return b, err
}

// FromEnvironmentVariables Открывает бакет, полностью обусловленный переменными окружения,
// в том числе и имя бакета тоже берется из MNML_BUCKET,
// MNML_KEY, MNML_SECRET - тут креденшалсы
func FromEnvironmentVariables() (*Bucket, error) {
	b, err := WithOptions(CredentialsFromEnv, NameFromEnv) // #todo а еще добавить в цепочку блять из файла
	return b, err
}

// FromKeys Открывает бакет, так что креденшалы задаются в открытом виде
func FromKeys(key, secret, session, name string) (*Bucket, error) {
	return nil, errors.New("функция FromFile Variables пока не воплощена")
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

// Опцефункции

type OptionsFunc func(b *Bucket) error

func NameFromEnv(b *Bucket) error {
	name, found := os.LookupEnv("MNML_BUCKET")
	if !found {
		fmt.Println("Не могу получить имя бакета из переменной окружения MNML_BUCKET")
	}

	return WithName(name)(b)
}

func WithName(name string) OptionsFunc {
	return func(b *Bucket) error {
		if name == "" {
			return nil // если имя пусток ничего не присваиваем дополнительно, может быть где-то в цепочке уже кто-то что-то ввел в имя
		}
		b.Name = name
		return nil
	}
}

func CredentialsFromEnv(b *Bucket) error {
	c, err := amazon.Client(amazon.LoadFromEnv)
	if err != nil {
		return err
	}
	b.s3client = c
	return nil
}

func ConfigFromFile(path string) OptionsFunc {
	return func(b *Bucket) error {
		return fmt.Errorf("не воплощено") //todo
	}
}

// #todo все эти функции выделить бы в отдельный файл с интерфейсом для AWS клиентировнных штук
