package main

import (
	"fmt"

	"github.com/thefrol/minimal/storage"
)

func main() {
	b, _ := storage.New("web-dir")
	err := b.UploadFile("../go.mod", "mnml-file")
	if err != nil {
		fmt.Println("error", err)
		return
	}
	r, err := b.GetString("mnml-file")
	if err != nil {
		fmt.Println(err)
		return
	}
	println(r)

	names, _ := b.Names()
	for _, s := range names {
		fmt.Println(s)
	}

}
