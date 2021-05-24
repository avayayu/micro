package dao

import (
	"reflect"
	"testing"
)

func Test_getTableFieldNameGormName(t *testing.T) {
	type args struct {
		model Model
	}
	tests := []struct {
		name string
		args args
		want map[string]string
	}{
		// TODO: Add test cases.

		{
			name: "test",
			args: args{
				model: &BluetoothDescribe{}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getTableFieldNameGormName(tt.args.model); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getTableFieldNameGormName() = %v, want %v", got, tt.want)
			}
		})
	}
}
