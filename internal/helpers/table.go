package helpers

import (
	"reflect"
	"strings"

	strcase "github.com/stoewer/go-strcase"
)

type Tabler interface {
	Table() string
}

func GetTable(m any) string {
	if m, ok := m.(Tabler); ok {
		rv := reflect.ValueOf(m)
		if rv.Kind() == reflect.Pointer && rv.IsNil() {
			rv = reflect.New(rv.Type().Elem())
			m = rv.Interface().(Tabler)
		}
		return m.Table()
	}
	t := reflect.TypeOf(m)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	name := strcase.SnakeCase(t.Name())
	if !strings.HasSuffix(name, "s") {
		name += "s"
	}
	return name
}

func GetTableSingular(m any) string {
	name := GetTable(m)

	if strings.HasSuffix(name, "s") {
		name = name[:len(name)-1]
	}
	return name
}
