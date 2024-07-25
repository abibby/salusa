package openapidoc

import (
	"encoding"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/abibby/salusa/internal/helpers"
	"github.com/go-openapi/spec"
)

var (
	// typeByteSlice             = reflect.TypeOf((*[]byte)(nil)).Elem()
	typeTime                  = reflect.TypeOf((*time.Time)(nil)).Elem()
	typeRune                  = reflect.TypeOf((*rune)(nil)).Elem()
	typeEncodingTextMarshaler = reflect.TypeOf((*encoding.TextMarshaler)(nil)).Elem()
	typeJSONMarshaler         = reflect.TypeOf((*json.Marshaler)(nil)).Elem()
	// typeModel                 = reflect.TypeOf((*model.Model)(nil)).Elem()
)

var formatMap = map[reflect.Type]string{
	typeEncodingTextMarshaler: "",
	typeTime:                  "date-time",
}

func RegisterFormat[T any](format string) {
	formatMap[reflect.TypeFor[T]()] = format
}

func Param(t reflect.Type, name, in string) (*spec.Parameter, error) {
	var err error
	param := &spec.Parameter{}
	param.Schema, err = Schema(t, false)
	if err != nil {
		return nil, err
	}
	param.Type = Type(t)
	param.Format = Format(t)
	param.Name = name
	param.In = in
	return param, nil
}
func Type(t reflect.Type) string {
	if t.Implements(typeEncodingTextMarshaler) {
		return "string"
	}
	switch t.Kind() {
	case reflect.Bool:
		return "boolean"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Uintptr:
		return "integer"
	case reflect.Float32, reflect.Float64:
		return "number"
	case reflect.Array, reflect.Slice:
		return "array"
	case reflect.Map, reflect.Struct:
		return "object"
	case reflect.String:
		return "string"
	default:
		return ""
	}
}

func Format(t reflect.Type) string {
	format, ok := formatMap[t]
	if ok {
		return format
	}

	switch t.Kind() {
	case reflect.Bool, reflect.Array, reflect.Slice, reflect.Map, reflect.Struct, reflect.String:
		return ""
	default:
		return t.Kind().String()
	}
}

func Schema(t reflect.Type, requireTag bool) (*spec.Schema, error) {

	switch t {
	case typeTime:
		return spec.DateTimeProperty(), nil
	case typeRune:
		return spec.CharProperty(), nil
	}

	if t.Implements(typeJSONMarshaler) {
		return nil, nil
	}
	if t.Implements(typeEncodingTextMarshaler) {
		return spec.StringProperty(), nil
	}
	// fmt.Printf("%v %v\n", t, t.Kind())
	switch t.Kind() {
	case reflect.Bool:
		return spec.BoolProperty(), nil
	case reflect.Int8:
		return spec.Int8Property(), nil
	case reflect.Int16, reflect.Uint8:
		return spec.Int16Property(), nil
	case reflect.Int32, reflect.Uint16:
		return spec.Int32Property(), nil
	case reflect.Int, reflect.Int64, reflect.Uint, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return spec.Int64Property(), nil
	case reflect.Float32:
		return spec.Float32Property(), nil
	case reflect.Float64:
		return spec.Float64Property(), nil
	case reflect.Array, reflect.Slice:
		item, err := Schema(t.Elem(), requireTag)
		if err != nil {
			return nil, err
		}
		return spec.ArrayProperty(item), nil
	case reflect.Map:
		item, err := Schema(t.Elem(), requireTag)
		if err != nil {
			return nil, err
		}
		if t.Key().Kind() != reflect.String {
			return nil, fmt.Errorf("map keys expected strings found %s", t.Key().Kind())
		}
		return spec.MapProperty(item), nil
	case reflect.Pointer:
		item, err := Schema(t.Elem(), requireTag)
		if err != nil {
			return nil, err
		}
		return item.AsNullable(), nil
	case reflect.String:
		return spec.StringProperty(), nil
	case reflect.Struct:

		schema := &spec.Schema{}
		schema.AddType("object", "")

		for _, field := range helpers.GetFields(t) {
			name, ok := field.Tag.Lookup("json")
			if ok {
				name = strings.Split(name, ",")[0]
				if name == "-" {
					continue
				}
			} else {
				if requireTag {
					continue
				}
				name = field.Name
			}
			item, err := Schema(field.Type, requireTag)
			if err != nil {
				return nil, err
			}
			schema.SetProperty(name, *item)
		}

		return schema, nil
	default:
		return nil, fmt.Errorf("unsupported type %v %v", t, t.Kind())
	}
}
