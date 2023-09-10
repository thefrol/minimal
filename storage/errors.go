package storage

import "fmt"

// KeyNotFound это тип ошибки, который возвращается если не найден нужный ключ в бакете, содержит информацию какой ключ и в каком бакете
type KeyNotFound struct {
	BucketName string
	Key        string
	Err        error
}

func (e *KeyNotFound) Error() string {
	return fmt.Sprintf("Ключ %v не найден в бакете %v", e.Key, e.BucketName)
}

func (e *KeyNotFound) Unwrap() error {
	return e.Err
}
