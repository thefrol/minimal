# `minimal` это минималистичный sdk для `Яндекс облака` на `Go`

Пакет призван упростить взаимодействие с облаком. В [официальной документации](https://cloud.yandex.ru/docs/storage/tools/aws-sdk-go) получение файла из бакета занимает 30 строк, в то время как с использованием `mnml/bucker`, это выглядит вот так. 

```go
b, _ := storage.New("my-bucket")
r, _ := b.GetString("my_file.txt")
fmt.Println(r)
```

Гораздо яснее, правда?

# Приступим

### Object Storage

Для начала подтребуется настройка aws, нужно ввести переменные окружения MNML_KEY, MNML_SECRET, а далее все просто:

```
import "github.com/thefrol/minimal/storage""

// Поключаемся к бакету по имени my-bucket
b, _ := storage.New("my-bucket")

//загружаем файл в бакет
b.UploadFile("test.txt", "test.txt")

//получаем файл из бакета, в виде строки. Для других типов подбробуйте функции вида Get...()
r, err := b.GetString("test.txt")
if err != nil {
	fmt.Println(err)
	return
}

В данном случае объект бакета работает как простой файл, и часто этого достаточно. 

```

## Настройка статических ключей(для AWS)

Получите статические ключи к сервисному аккаунту и заполните переменные окружения MNML_KEY, MNNL_SECRET

### Бакеты

# Автор 

[Дмитрий Фроленко](https://github.com/thefrol) 2023
