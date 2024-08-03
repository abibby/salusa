package jsoncolumn

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"reflect"
)

func Scan(dst, src any) error {
	dstType := reflect.TypeOf(dst)
	if dstType.Kind() != reflect.Pointer {
		return &json.InvalidUnmarshalError{Type: dstType}
	}
	switch src := src.(type) {
	case []byte:
		return json.Unmarshal(src, dst)
	case nil:
		reflect.ValueOf(dst).Elem().Set(reflect.Zero(dstType.Elem()))
		return nil
	default:
		return fmt.Errorf("invalid type %s", reflect.TypeOf(src))
	}
}
func Value(v any) (driver.Value, error) {
	return json.Marshal(v)
}
