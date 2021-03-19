package collection

import (
	"fmt"
	"testing"
)

type test struct {
	HAHA  string
	HAHA2 string
}

type test2 struct {
	HAHA  string
	HAHA2 string
	HAHA3 *test
}

func TestGetStructColumFromSlice(t *testing.T) {
	type args struct {
		in     interface{}
		column string
		out    interface{}
	}
	var out []string
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name:    "test",
			wantErr: false,
			args: args{
				in: []test{
					{
						HAHA:  "A",
						HAHA2: "A2",
					},
					{
						HAHA:  "C",
						HAHA2: "C2",
					},
				},
				column: "HAHA",
				out:    &out,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := GetStructColumFromSlice(tt.args.in, tt.args.column, tt.args.out); (err != nil) != tt.wantErr {
				t.Errorf("GetStructColumFromSlice() error = %v, wantErr %v", err, tt.wantErr)
			}

			fmt.Println(tt.args.out)
		})
	}
}
