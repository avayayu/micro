package models

import (
	"reflect"
	"sort"
	"testing"
)

func TestDiff(t *testing.T) {

	a := []Int64Str{1, 2, 3, 5, 7}
	b := []Int64Str{2, 4, 10}
	a2 := []Int64Str{55555, 2, 3, 10000}
	b2 := []Int64Str{10000, 3, 2, 1000000, 3123123, 55555, 21231}
	want := IDLIST{1000000, 3123123, 21231}
	sort.Sort(&(want))
	type args struct {
		a IDLIST
		b IDLIST
	}
	tests := []struct {
		name string
		args args
		want IDLIST
	}{
		// TODO: Add test cases.
		{
			name: "test",
			args: args{
				a: a,
				b: b,
			},
			want: []Int64Str{4, 10},
		},
		{
			name: "test2",
			args: args{
				a: a2,
				b: b2,
			},
			want: want,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Diff(tt.args.a, tt.args.b); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Diff() = %v, want %v", got, tt.want)
			}
		})
	}
}
