package storage

var (
	// список оберток для ошибок s3
	KeyNotFound = wrapError("Ключ не найден в бакете")
)

// StorageError оперделяет какая ошибка произошла, и какая ошибка в s3 её вызвала
type StorageError struct {
	Text string
	Err  error
}

func (se StorageError) Error() string {
	return se.Text
}

// wrapError Оборачивает ошибку оригинальную, ошибкой из этого пакета. Так сохраняется преемственность ошибок от S3.
// newtext опеределяет новый текст ошибки
func wrapError(newtext string) func(error) error {
	return func(parent error) error {
		return StorageError{
			Text: newtext,
			Err:  parent,
		}
	}
}
