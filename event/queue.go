package event

import (
	"context"
	"fmt"
	"reflect"

	"github.com/abibby/salusa/di"
	"github.com/abibby/salusa/kernel"
)

type Queue interface {
	Push(e Event) error
	Pop(events map[EventType]reflect.Type) (Event, error)
}

type QueueConfiger interface {
	QueueConfig() Config
}
type Config interface {
	Queue() Queue
}

func Register[T kernel.KernelConfig](ctx context.Context) error {
	di.RegisterLazySingletonWith(ctx, func(cfg T) (Queue, error) {
		var cfgAny any = cfg
		cfger, ok := cfgAny.(QueueConfiger)
		if !ok {
			return nil, fmt.Errorf("config not instance of event.QueueConfiger")
		}
		return cfger.QueueConfig().Queue(), nil
	})
	return nil
}
