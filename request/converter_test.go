package request

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConverters(t *testing.T) {
	tests := []struct {
		name      string
		conv      Converter
		value     string
		want      any
		wantValid bool
	}{
		{"bool-on", convertBool, "on", true, true},
		{"bool-true", convertBool, "true", true, true},
		{"bool-1", convertBool, "1", true, true},
		{"bool-false", convertBool, "false", false, true},
		{"bool-invalid", convertBool, "not bool", nil, false},
		{"float32", convertFloat32, "1.5", float32(1.5), true},
		{"float32-invalid", convertFloat32, "not float", nil, false},
		{"float64", convertFloat64, "1.5", 1.5, true},
		{"float64-invalid", convertFloat64, "not float", nil, false},
		{"int", convertInt, "15", 15, true},
		{"int-invalid", convertInt, "not int", nil, false},
		{"int8", convertInt8, "15", int8(15), true},
		{"int8-invalid", convertInt8, "not int", nil, false},
		{"int16", convertInt16, "15", int16(15), true},
		{"int16-invalid", convertInt16, "not int", nil, false},
		{"int32", convertInt32, "15", int32(15), true},
		{"int32-invalid", convertInt32, "not int", nil, false},
		{"int64", convertInt64, "15", int64(15), true},
		{"int64-invalid", convertInt64, "not int", nil, false},
		{"string", convertString, "foo", "foo", true},
		{"uint", convertUint, "15", uint(15), true},
		{"uint-invalid", convertUint, "not uint", nil, false},
		{"uint8", convertUint8, "15", uint8(15), true},
		{"uint8-invalid", convertUint8, "not uint", nil, false},
		{"uint16", convertUint16, "15", uint16(15), true},
		{"uint16-invalid", convertUint16, "not uint", nil, false},
		{"uint32", convertUint32, "15", uint32(15), true},
		{"uint32-invalid", convertUint32, "not uint", nil, false},
		{"uint64", convertUint64, "15", uint64(15), true},
		{"uint64-invalid", convertUint64, "not uint", nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.conv(tt.value)
			if tt.wantValid {
				assert.True(t, got.IsValid())
				assert.Equal(t, reflect.ValueOf(tt.want).Interface(), got.Interface())
			} else {
				assert.Equal(t, invalidValue, got)
			}
		})
	}
}

func TestDecodeBuiltin(t *testing.T) {
	t.Run("unsupported kind", func(t *testing.T) {
		type S struct{}
		_, err := decode(reflect.TypeFor[S](), []string{"foo"})
		assert.Error(t, err)
	})

	t.Run("int8 via decode", func(t *testing.T) {
		v, err := decode(reflect.TypeFor[int8](), []string{"5"})
		assert.NoError(t, err)
		assert.Equal(t, int8(5), v.Interface())
	})
}
