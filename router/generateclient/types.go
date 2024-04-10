package generateclient

import (
	"fmt"
	"io"
	"reflect"
)

func toTsType(w io.Writer, t reflect.Type, overrides map[reflect.Type]string) error {

	if overrides != nil {
		if o, ok := overrides[t]; ok {
			_, err := w.Write([]byte(o))
			return err
		}
	}

	switch t.Kind() {
	case reflect.Bool:
		_, err := fmt.Fprint(w, "boolean")
		return err
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Uintptr,
		reflect.Float32, reflect.Float64:
		_, err := fmt.Fprint(w, "number")
		return err
	case reflect.Map:
		return fmt.Errorf("map not supported")
	case reflect.Pointer:
		if t.Elem().Kind() == reflect.Struct {
			return toTsTypeStruct(w, t.Elem(), "json", overrides)
		}

		err := toTsType(w, t.Elem(), overrides)
		if err != nil {
			return err
		}
		_, err = fmt.Fprint(w, " | null")
		return err
	case reflect.Slice:
		_, err := fmt.Fprint(w, "Array<")
		if err != nil {
			return err
		}
		err = toTsType(w, t.Elem(), overrides)
		if err != nil {
			return err
		}
		_, err = fmt.Fprint(w, ">")
		return err
	case reflect.String:
		_, err := fmt.Fprint(w, "string")
		return err
	case reflect.Struct:
		return toTsTypeStruct(w, t, "json", overrides)
	default:
		return fmt.Errorf("unsupported type %s", t)
	}
}

func toTsTypeStruct(w io.Writer, t reflect.Type, tag string, overrides map[reflect.Type]string) error {
	_, err := w.Write([]byte{'{'})
	if err != nil {
		return err
	}
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		name, ok := f.Tag.Lookup(tag)
		if !ok {
			continue
		}
		if i > 0 {
			_, err = w.Write([]byte(","))
			if err != nil {
				return err
			}
		}
		_, err = w.Write([]byte(name + ":"))
		if err != nil {
			return err
		}
		err = toTsType(w, f.Type, overrides)
		if err != nil {
			return err
		}

	}
	_, err = w.Write([]byte{'}'})
	if err != nil {
		return err
	}
	return nil
}
