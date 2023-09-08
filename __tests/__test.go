package main

import (
	"fmt"

	"github.com/thefrol/minimal/storage"
)

func main() {
	b, _ := storage.New("web-dir")
	b.UploadFile("test.go", "mnml-file")
	r, err := b.GetString("123")
	if err != nil {
		fmt.Println(err)
		return
	}
	println(r)

}
