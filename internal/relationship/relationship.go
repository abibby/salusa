package relationship

import (
	"reflect"
)

type Relationship interface {
	Initialize(self any, field reflect.StructField) error
	Loaded() bool
}
