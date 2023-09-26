package event

import "reflect"

type Queue interface {
	Push(e Event) error
	Pop(events map[EventType]reflect.Type) (Event, error)
}
