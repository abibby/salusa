package request

import (
	"encoding"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/abibby/salusa/database/model"
	"github.com/abibby/salusa/internal/helpers"
	"github.com/abibby/salusa/openapidoc"
	"github.com/go-openapi/spec"
)

var (
	// typeByteSlice             = reflect.TypeOf((*[]byte)(nil)).Elem()
	typeTime                  = reflect.TypeOf((*time.Time)(nil)).Elem()
	typeRune                  = reflect.TypeOf((*rune)(nil)).Elem()
	typeEncodingTextMarshaler = reflect.TypeOf((*encoding.TextMarshaler)(nil)).Elem()
	typeJSONMarshaler         = reflect.TypeOf((*json.Marshaler)(nil)).Elem()
	typeModel                 = reflect.TypeOf((*model.Model)(nil)).Elem()
)

var _ openapidoc.Operationer = (*RequestHandler[any, any])(nil)

// Operation implements openapidoc.Operationer.
func (h *RequestHandler[TRequest, TResponse]) Operation() (*spec.Operation, error) {
	var err error
	op := spec.NewOperation("")
	// op.Parameters
	// spec.NewResponse()

	op.Responses = &spec.Responses{}
	op.Parameters, err = newAPIRequest[TRequest]()
	if err != nil {
		return nil, err
	}
	op.Responses.Default, err = newAPIResponse[TResponse]()
	if err != nil {
		return nil, err
	}
	return op, nil
}

func newAPIRequest[T any]() ([]spec.Parameter, error) {
	params := []spec.Parameter{}
	var emptyRequest *T
	t := reflect.TypeOf(emptyRequest).Elem()

	schema, err := buildSchema(t, true)
	if err != nil {
		return nil, err
	}
	if len(schema.Properties) > 0 {
		params = append(params, *spec.BodyParam("Body", schema))
	}
	for _, field := range helpers.GetFields(t) {
		if name, ok := field.Tag.Lookup("query"); ok {
			param, err := buildParam(field.Type, name, "query")
			if err != nil {
				return nil, err
			}
			params = append(params, *param)
		}

		if name, ok := field.Tag.Lookup("path"); ok {
			param, err := buildParam(field.Type, name, "path")
			if err != nil {
				return nil, err
			}
			params = append(params, *param)
		}

		if field.Type.Implements(typeModel) {
			if name, ok := field.Tag.Lookup("inject"); ok && name != "" {
				field.Type.FieldByName(name)
				param, err := buildParam(field.Type, name, "path")
				if err != nil {
					return nil, err
				}
				params = append(params, *param)
			}
		}
	}

	return params, nil
}
func newAPIResponse[T any]() (*spec.Response, error) {
	resp := spec.NewResponse()
	emptyResponse, err := helpers.NewOf[T]()
	if err != nil {
		return nil, nil
	}

	resp.Schema, err = buildSchema(reflect.TypeOf(emptyResponse), false)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func buildParam(t reflect.Type, name, in string) (*spec.Parameter, error) {
	var err error
	param := &spec.Parameter{}
	param.Schema, err = buildSchema(t, false)
	if err != nil {
		return nil, err
	}
	param.Type = buildType(t)
	param.Format = buildFormat(t)
	param.Name = name
	param.In = in
	return param, nil
}
func buildType(t reflect.Type) string {
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

func buildFormat(t reflect.Type) string {
	if t.Implements(typeEncodingTextMarshaler) {
		return ""
	}
	switch t.Kind() {
	case reflect.Bool, reflect.Array, reflect.Slice, reflect.Map, reflect.Struct, reflect.String:
		return ""
	default:
		return t.Kind().String()
	}
}

func buildSchema(t reflect.Type, requireTag bool) (*spec.Schema, error) {

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
		item, err := buildSchema(t.Elem(), requireTag)
		if err != nil {
			return nil, err
		}
		return spec.ArrayProperty(item), nil
	case reflect.Map:
		item, err := buildSchema(t.Elem(), requireTag)
		if err != nil {
			return nil, err
		}
		if t.Key().Kind() != reflect.String {
			return nil, fmt.Errorf("map keys expected strings found %s", t.Key().Kind())
		}
		return spec.MapProperty(item), nil
	case reflect.Pointer:
		item, err := buildSchema(t.Elem(), requireTag)
		if err != nil {
			return nil, err
		}
		return item.AsNullable(), nil
	case reflect.String:
		return spec.StringProperty(), nil
	case reflect.Struct:

		schema := &spec.Schema{}
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
			item, err := buildSchema(field.Type, requireTag)
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
