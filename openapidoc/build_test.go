package openapidoc

import (
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/go-openapi/spec"
	"github.com/stretchr/testify/assert"
)

type customTextMarshaler struct{}

func (customTextMarshaler) MarshalText() ([]byte, error) { return []byte("x"), nil }

type customJSONMarshaler struct{}

func (customJSONMarshaler) MarshalJSON() ([]byte, error) { return []byte("{}"), nil }

type testDocStruct struct {
	Name  string `json:"name"`
	Skip  string `json:"-"`
	NoTag string
}

func TestType(t *testing.T) {
	tests := []struct {
		name string
		typ  any
		want string
	}{
		{"text marshaler", customTextMarshaler{}, "string"},
		{"bool", true, "boolean"},
		{"int", int(1), "integer"},
		{"int8", int8(1), "integer"},
		{"int16", int16(1), "integer"},
		{"int32", int32(1), "integer"},
		{"int64", int64(1), "integer"},
		{"uint", uint(1), "integer"},
		{"uint8", uint8(1), "integer"},
		{"uint16", uint16(1), "integer"},
		{"uint32", uint32(1), "integer"},
		{"uint64", uint64(1), "integer"},
		{"uintptr", uintptr(1), "integer"},
		{"float32", float32(1), "number"},
		{"float64", float64(1), "number"},
		{"array", [2]int{}, "array"},
		{"slice", []int{}, "array"},
		{"map", map[string]int{}, "object"},
		{"struct", struct{}{}, "object"},
		{"string", "x", "string"},
		{"pointer", new(int), "integer"},
		{"chan", make(chan int), ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, Type(reflect.TypeOf(tt.typ)))
		})
	}
}

func TestFormat(t *testing.T) {
	tests := []struct {
		name string
		typ  reflect.Type
		want string
	}{
		{"time", typeTime, "date-time"},
		{"time pointer", reflect.TypeFor[*time.Time](), ""},
		{"int32 pointer", reflect.TypeFor[*int32](), "int32"},
		{"text marshaler", reflect.TypeFor[customTextMarshaler](), ""},
		{"bool", reflect.TypeFor[bool](), ""},
		{"int", reflect.TypeFor[int](), ""},
		{"string", reflect.TypeFor[string](), ""},
		{"int32", reflect.TypeFor[int32](), "int32"},
		{"float64", reflect.TypeFor[float64](), "float64"},
		{"chan", reflect.TypeFor[chan int](), "chan"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, Format(tt.typ))
		})
	}
}

func TestRegisterFormat(t *testing.T) {
	type customFormatType struct{}
	RegisterFormat[customFormatType]("custom-format")
	assert.Equal(t, "custom-format", Format(reflect.TypeFor[customFormatType]()))
}

func TestSchema(t *testing.T) {
	t.Run("registered types", func(t *testing.T) {
		s, err := Schema(reflect.TypeFor[time.Time](), false)
		assert.NoError(t, err)
		assert.NotNil(t, s)

		s, err = Schema(reflect.TypeFor[rune](), false)
		assert.NoError(t, err)
		assert.NotNil(t, s)
	})

	t.Run("json marshaler", func(t *testing.T) {
		s, err := Schema(reflect.TypeFor[customJSONMarshaler](), false)
		assert.NoError(t, err)
		assert.Equal(t, "unknown", s.Type[0])
	})

	t.Run("text marshaler", func(t *testing.T) {
		s, err := Schema(reflect.TypeFor[customTextMarshaler](), false)
		assert.NoError(t, err)
		assert.Equal(t, "string", s.Type[0])
	})

	t.Run("bool", func(t *testing.T) {
		s, err := Schema(reflect.TypeFor[bool](), false)
		assert.NoError(t, err)
		assert.Equal(t, "boolean", s.Type[0])
	})

	t.Run("integers", func(t *testing.T) {
		tests := []struct {
			typ    reflect.Type
			typStr string
			format string
		}{
			{reflect.TypeFor[int8](), "integer", "int8"},
			{reflect.TypeFor[int16](), "integer", "int16"},
			{reflect.TypeFor[uint8](), "integer", "int16"},
			{reflect.TypeFor[int32](), "string", ""},
			{reflect.TypeFor[uint16](), "integer", "int32"},
			{reflect.TypeFor[int](), "integer", "int64"},
			{reflect.TypeFor[int64](), "integer", "int64"},
			{reflect.TypeFor[uint](), "integer", "int64"},
			{reflect.TypeFor[uint32](), "integer", "int64"},
			{reflect.TypeFor[uint64](), "integer", "int64"},
			{reflect.TypeFor[uintptr](), "integer", "int64"},
		}
		for _, tt := range tests {
			s, err := Schema(tt.typ, false)
			assert.NoError(t, err)
			assert.Equal(t, tt.typStr, s.Type[0])
			assert.Equal(t, tt.format, s.Format)
		}
	})

	t.Run("floats", func(t *testing.T) {
		s, err := Schema(reflect.TypeFor[float32](), false)
		assert.NoError(t, err)
		assert.Equal(t, "number", s.Type[0])
		assert.Equal(t, "float", s.Format)

		s, err = Schema(reflect.TypeFor[float64](), false)
		assert.NoError(t, err)
		assert.Equal(t, "number", s.Type[0])
		assert.Equal(t, "double", s.Format)
	})

	t.Run("string", func(t *testing.T) {
		s, err := Schema(reflect.TypeFor[string](), false)
		assert.NoError(t, err)
		assert.Equal(t, "string", s.Type[0])
	})

	t.Run("slice", func(t *testing.T) {
		s, err := Schema(reflect.TypeFor[[]string](), false)
		assert.NoError(t, err)
		assert.Equal(t, "array", s.Type[0])
		assert.Equal(t, "string", s.Items.Schema.Type[0])
	})

	t.Run("map string keys", func(t *testing.T) {
		s, err := Schema(reflect.TypeFor[map[string]int](), false)
		assert.NoError(t, err)
		assert.Equal(t, "object", s.Type[0])
	})

	t.Run("map non-string keys", func(t *testing.T) {
		_, err := Schema(reflect.TypeFor[map[int]string](), false)
		assert.Error(t, err)
	})

	t.Run("pointer", func(t *testing.T) {
		s, err := Schema(reflect.TypeFor[*string](), false)
		assert.NoError(t, err)
		assert.True(t, s.Nullable)
		assert.Equal(t, "string", s.Type[0])
	})

	t.Run("struct", func(t *testing.T) {
		s, err := Schema(reflect.TypeFor[testDocStruct](), false)
		assert.NoError(t, err)
		_, ok := s.Properties["name"]
		assert.True(t, ok)
		_, ok = s.Properties["skip"]
		assert.False(t, ok)
		_, ok = s.Properties["NoTag"]
		assert.True(t, ok)
	})

	t.Run("struct requires tag", func(t *testing.T) {
		s, err := Schema(reflect.TypeFor[testDocStruct](), true)
		assert.NoError(t, err)
		_, ok := s.Properties["name"]
		assert.True(t, ok)
		_, ok = s.Properties["NoTag"]
		assert.False(t, ok)
	})

	t.Run("unsupported type", func(t *testing.T) {
		_, err := Schema(reflect.TypeFor[chan int](), false)
		assert.Error(t, err)
	})

	t.Run("slice of unsupported", func(t *testing.T) {
		_, err := Schema(reflect.TypeFor[[]chan int](), false)
		assert.Error(t, err)
	})

	t.Run("map of unsupported", func(t *testing.T) {
		_, err := Schema(reflect.TypeFor[map[string]chan int](), false)
		assert.Error(t, err)
	})

	t.Run("pointer to unsupported", func(t *testing.T) {
		_, err := Schema(reflect.TypeFor[*chan int](), false)
		assert.Error(t, err)
	})

	t.Run("struct field unsupported", func(t *testing.T) {
		_, err := Schema(reflect.TypeFor[struct {
			C chan int `json:"c"`
		}](), false)
		assert.Error(t, err)
	})
}

func TestRegisterSchema(t *testing.T) {
	type customSchemaType struct{}
	RegisterSchema[customSchemaType](spec.StringProperty())
	s, err := Schema(reflect.TypeFor[customSchemaType](), false)
	assert.NoError(t, err)
	assert.Equal(t, "string", s.Type[0])
}

func TestResponse(t *testing.T) {
	t.Run("registered response", func(t *testing.T) {
		r, err := Response(reflect.TypeFor[*http.Response]())
		assert.NoError(t, err)
		assert.NotNil(t, r)
	})

	t.Run("custom registered response", func(t *testing.T) {
		type customResp struct{}
		RegisterResponse[customResp](spec.NewResponse().WithDescription("custom"))
		r, err := Response(reflect.TypeFor[customResp]())
		assert.NoError(t, err)
		assert.Equal(t, "custom", r.Description)
	})

	t.Run("built from schema", func(t *testing.T) {
		r, err := Response(reflect.TypeFor[testDocStruct]())
		assert.NoError(t, err)
		assert.NotNil(t, r)
		assert.NotNil(t, r.Schema)
	})

	t.Run("schema error logs warn", func(t *testing.T) {
		r, err := Response(reflect.TypeFor[chan int]())
		assert.NoError(t, err)
		assert.NotNil(t, r)
		assert.Nil(t, r.Schema)
	})
}

func TestRegisterContentType(t *testing.T) {
	_, ok := GetContentType[testDocStruct]()
	assert.False(t, ok)

	RegisterContentType[testDocStruct]("application/json")
	ct, ok := GetContentType[testDocStruct]()
	assert.True(t, ok)
	assert.Equal(t, "application/json", ct)
}

func TestParam(t *testing.T) {
	p, err := Param(reflect.TypeFor[string](), "foo", "query")
	assert.NoError(t, err)
	assert.Equal(t, "foo", p.Name)
	assert.Equal(t, "query", p.In)
	assert.Equal(t, "string", p.Type)

	_, err = Param(reflect.TypeFor[chan int](), "bad", "query")
	assert.Error(t, err)
}
