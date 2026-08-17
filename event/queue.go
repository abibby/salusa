package event

import (
	"context"
	"reflect"
)

type Queue interface {
	Push(ctx context.Context, e Event) error
	Pop(ctx context.Context, events map[EventType]reflect.Type) (Event, error)
}
