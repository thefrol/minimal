package bucket

import (
	"os"
	"reflect"
	"strings"
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

func TestBucket_GetString(t *testing.T) {
	type args struct {
		objectKey string
		content   string
	}
	tests := []struct {
		name          string
		createdFile   args
		gettedFile    args
		wantErr       bool
		concreteError error // if nil any error is good
	}{
		{
			name:          "simple",
			createdFile:   args{objectKey: testFile, content: uploadContent},
			gettedFile:    args{objectKey: testFile, content: uploadContent},
			wantErr:       false,
			concreteError: nil,
		},
		{
			name:          "no key, specific error #1",
			createdFile:   args{objectKey: testFile, content: uploadContent},
			gettedFile:    args{objectKey: anotherFile, content: uploadContent},
			wantErr:       true,
			concreteError: &KeyNotFound{},
		},
		{
			name:          "bad key key, any error #2",
			createdFile:   args{objectKey: testFile, content: uploadContent},
			gettedFile:    args{objectKey: "///wtf", content: uploadContent},
			wantErr:       true,
			concreteError: nil,
		},
		{
			name:          "bad key key, any error #3",
			createdFile:   args{objectKey: testFile, content: uploadContent},
			gettedFile:    args{objectKey: "/...#$%^&*wtf", content: uploadContent},
			wantErr:       true,
			concreteError: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := cleanBucket()
			require.NoError(t, err)

			err = b.Put(strings.NewReader(tt.createdFile.content), tt.createdFile.objectKey)
			require.NoError(t, err)

			s, err := b.GetString(tt.gettedFile.objectKey)
			if tt.wantErr {
				assert.Error(t, err)
				if tt.concreteError != nil {
					assert.ErrorAs(t, err, &tt.concreteError)
				}
				return

			}
			require.NoError(t, err)

			assert.Equal(t, tt.createdFile.content, s, "Got wrong content")

		})
	}
}

func TestBucket_Deletefile(t *testing.T) {

	tests := []struct {
		name         string
		created      []string
		keysToDelete []string
		left         []string
		wantErr      bool
	}{
		{
			name:         "simple",
			created:      []string{testFile},
			keysToDelete: []string{testFile},
			left:         []string{},
			wantErr:      false,
		},
		{
			name:         "not existent key, no error",
			created:      []string{testFile},
			keysToDelete: []string{anotherFile},
			left:         []string{testFile},
			wantErr:      false,
		},
		{
			name:         "multiple",
			created:      []string{testFile, anotherFile, "third.file"},
			keysToDelete: []string{anotherFile, testFile},
			left:         []string{"third.file"},
			wantErr:      false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := cleanBucket()
			require.NoError(t, err)

			for _, f := range tt.created {
				err := b.Put(strings.NewReader(uploadContent), f)
				require.NoError(t, err)

			}

			wasError := false
			for _, f := range tt.keysToDelete {
				err := b.Delete(f)
				if tt.wantErr {
					wasError = true
				} else {
					require.NoError(t, err)
				}
			}

			if tt.wantErr {
				assert.True(t, tt.wantErr, wasError, "We expected to have an error here")
			}
			actual, err := b.Names()
			require.NoError(t, err)

			assert.Truef(t, reflect.DeepEqual(tt.left, actual), "Files left in bucket should be %v, but its %v", tt.left, actual)
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

// #todo все эти функции только бакета не существует
