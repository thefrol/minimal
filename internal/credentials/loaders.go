package credentials

import (
	"os"

	"gopkg.in/yaml.v3"
)

//FromFile загружает секреты из файла
func FromFile(path string) (*Credentials, error) {

	r, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	d := yaml.NewDecoder(r)
	c := Credentials{}
	err = d.Decode(&c)
	if err != nil {
		return nil, err
	}

	return &c, nil
}
