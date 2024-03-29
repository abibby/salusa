package helpers

import (
	"fmt"
	"reflect"
)

type PrimaryKeyer interface {
	PrimaryKey() []string
}

func PrimaryKey(m any) []string {
	if m, ok := m.(PrimaryKeyer); ok {
		return m.PrimaryKey()
	}

	t := reflect.TypeOf(m)
	primary, fallback := primaryKey(t)
	if len(primary) == 0 {
		return []string{fallback}
	}
	return primary
}

func PrimaryKeyValue(m any) ([]any, error) {
	modelValue := reflect.ValueOf(m)
	keys := PrimaryKey(m)
	values := make([]any, len(keys))
	for i, key := range keys {
		v, err := RGetValue(modelValue, key)
		if err != nil {
			return nil, fmt.Errorf("can't find key %s on %s: %w", key, modelValue.Type(), err)
		}
		values[i] = v.Interface()
	}
	return values, nil
}

func primaryKey(t reflect.Type) ([]string, string) {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	fallback := ""
	primary := []string{}
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if f.Anonymous {
			p, f := primaryKey(f.Type)
			if fallback == "" {
				fallback = f
			}
			primary = append(primary, p...)
			continue
		}
		if !f.IsExported() {
			continue
		}
		tag := DBTag(f)
		if fallback == "" {
			fallback = tag.Name
		}
		if tag.Primary {
			primary = append(primary, tag.Name)
		}
	}

	return primary, fallback
}

func Includes[T comparable](arr []T, v T) bool {
	for _, e := range arr {
		if v == e {
			return true
		}
	}
	return false
}
