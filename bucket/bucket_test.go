package bucket

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestBucket_UploadFile(t *testing.T) {
	type args struct {
		fileName  string
		objectKey string
	}
	type createdFile struct {
		fileName string
		content  string
	}
	tests := []struct {
		name        string
		createdFile createdFile
		args        args
		wantErr     bool
	}{
		{
			name:        "simple",
			createdFile: createdFile{fileName: testFile, content: uploadContent},
			args:        args{fileName: testFile, objectKey: testFile},
			wantErr:     false,
		},
		{
			name:        "dirs",
			createdFile: createdFile{fileName: testFile, content: uploadContent},
			args:        args{fileName: testFile, objectKey: testFile},
			wantErr:     false,
		},
		{
			name:        "no file",
			createdFile: createdFile{fileName: testFile, content: uploadContent},
			args:        args{fileName: anotherFile, objectKey: testFile},
			wantErr:     true,
		},
		// а попробовать загрузить несуществующий файл?
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := cleanBucket()
			require.NoError(t, err)

			f, err := os.Create(tt.createdFile.fileName)
			require.NoError(t, err)
			defer os.Remove(tt.createdFile.fileName)

			_, err = f.Write([]byte(tt.createdFile.content))
			require.NoError(t, err)
			defer f.Close()

			err = b.UploadFile(tt.args.fileName, tt.args.objectKey)
			if tt.wantErr {
				require.Error(t, err)
				return
			} else {
				require.NoError(t, err)
			}

			s, err := b.GetString(tt.args.objectKey)

			require.NoError(t, err)
			assert.Equal(t, uploadContent, s, "Got wrong content")

		})
	}
}

const uploadContent = "This is a test content"
const (
	testFile    = "test.txt"
	anotherFile = "another.txt"
)

const credentialsPath = ".test_credentials.bucket"

var b *Bucket

func cleanBucket() (*Bucket, error) {
	// а вообще мы моглибы и создавать бакет, может такая функция даже будет
	c := struct {
		Key     string `yaml:"key"`
		Secret  string `yaml:"secret"`
		Session string `yaml:"session"`
		Bucket  string `yaml:"bucket"`
	}{}

	r, err := os.Open(credentialsPath)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	d := yaml.NewDecoder(r)
	err = d.Decode(&c)
	if err != nil {
		return nil, err
	}

	b, err = FromKeys(c.Key, c.Secret, c.Session, c.Bucket)

	if err != nil {
		return nil, err
	}

	files, err := b.Names()
	if err != nil {
		return nil, err
	}

	for _, fl := range files {
		err := b.Delete(fl)
		if err != nil {
			return nil, err
		}
	}

	return b, nil
}
