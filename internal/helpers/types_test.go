package helpers

import (
	"reflect"
	"testing"
)

type Foo struct{}

type FooPtr *Foo

func TestRNewOf(t *testing.T) {
	type args struct {
		t reflect.Type
	}
	tests := []struct {
		name    string
		args    args
		want    any
		wantErr bool
	}{
		{
			name: "a",
			args: args{reflect.TypeFor[FooPtr]()},
			want: FooPtr(&Foo{}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := RNewOf(tt.args.t)
			if (err != nil) != tt.wantErr {
				t.Errorf("RNewOf() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.Interface(), tt.want) {
				t.Errorf("RNewOf() = %v, want %v", got, tt.want)
			}
		})
	}
}
