package helpers

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestNewOf(t *testing.T) {
	t.Run("pointer", func(t *testing.T) {
		v, err := NewOf[*Foo]()
		require.NoError(t, err)
		assert.IsType(t, &Foo{}, v)
	})

	t.Run("struct", func(t *testing.T) {
		v, err := NewOf[Foo]()
		require.NoError(t, err)
		assert.IsType(t, Foo{}, v)
	})

	t.Run("interface", func(t *testing.T) {
		_, err := NewOf[any]()
		assert.Error(t, err)
	})

	t.Run("int", func(t *testing.T) {
		v, err := NewOf[int]()
		require.NoError(t, err)
		assert.Equal(t, 0, v)
	})
}

func TestCreateFor(t *testing.T) {
	t.Run("pointer", func(t *testing.T) {
		v := CreateFor[*Foo]()
		assert.Equal(t, reflect.Pointer, v.Kind())
	})

	t.Run("struct", func(t *testing.T) {
		v := CreateFor[Foo]()
		assert.Equal(t, reflect.Struct, v.Kind())
	})
}

func TestCreate(t *testing.T) {
	t.Run("pointer", func(t *testing.T) {
		v := Create(reflect.TypeFor[*Foo]())
		assert.Equal(t, reflect.Pointer, v.Kind())
		assert.IsType(t, &Foo{}, v.Interface())
	})

	t.Run("struct", func(t *testing.T) {
		v := Create(reflect.TypeFor[Foo]())
		assert.Equal(t, reflect.Struct, v.Kind())
	})
}

func TestZero(t *testing.T) {
	assert.Equal(t, 0, Zero[int]())
	assert.Equal(t, "", Zero[string]())
	assert.Nil(t, Zero[*Foo]())
}
