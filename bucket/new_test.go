package bucket

// так же добавить тест который смотритчто если не созданы кредентиалы, то он может взять из домашней директории .aws
// домашнюю дерикторию сделаить липовую через переменные окружения

import (
	"testing"
)

func Test_NewNaming(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		//wantB   *Bucket
		wantErr bool
	}{
		{
			name: "empty name",
			args: args{name: ""},
			//wantB:   nil,
			wantErr: true,
		},
		{
			name: "positive name",
			args: args{name: "buck"},
			//wantB:   nil,
			wantErr: false,
		},
		{
			name: "positive name",
			args: args{name: "buck.2"},
			//wantB:   nil,
			wantErr: false,
		},
		{
			name: "negative name #2 starts with digit",
			args: args{name: "2buck-2"},
			//wantB:   nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := New(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// if !reflect.DeepEqual(gotB, tt.wantB) {
			// 	t.Errorf("New() = %v, want %v", gotB, tt.wantB)
			// }
		})
	}
}
