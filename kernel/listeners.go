package kernel

import (
	"context"
	"reflect"

	"log/slog"

	"github.com/abibby/salusa/di"
	"github.com/abibby/salusa/event"
)

type JobHandler[E event.Event] interface {
	Handle(ctx context.Context, event E) error
}

type runner interface {
	UpdateValue(v event.Event) bool
	Run(ctx context.Context, dp *di.DependencyProvider) error
	EventType() reflect.Type
}
type Listener struct {
	eventType event.EventType
	runner    runner
}

type job[E event.Event] struct {
	value       E
	handlerType reflect.Type
}

func (j *job[E]) Run(ctx context.Context, dp *di.DependencyProvider) error {
	t := j.handlerType
	var h JobHandler[E]
	if t.Kind() == reflect.Pointer {
		h = reflect.New(t.Elem()).Interface().(JobHandler[E])
	} else {
		h = reflect.New(t).Elem().Interface().(JobHandler[E])
	}

	err := dp.Fill(ctx, h)
	if err != nil {
		return err
	}
	return h.Handle(ctx, j.value)
}

func (j *job[E]) UpdateValue(v event.Event) bool {
	ev, ok := v.(E)
	if !ok {
		return false
	}
	j.value = ev
	return true
}
func (j *job[E]) EventType() reflect.Type {
	var e E
	return reflect.TypeOf(e)
}

func NewListener[H JobHandler[E], E event.Event]() *Listener {
	var e E
	var h H

	return &Listener{
		eventType: e.Type(),
		runner: &job[E]{
			value:       e,
			handlerType: reflect.TypeOf(h),
		},
	}
}

func (k *Kernel) RunListeners(ctx context.Context) {
	events := map[event.EventType]reflect.Type{}
	for eventType, runners := range k.listeners {
		events[eventType] = runners[0].EventType()
	}

	for {
		e, err := k.queue.Pop(events)
		if err != nil {
			slog.Warn("could not pop event off queue", slog.Any("error", err))
			continue
		}
		runners, ok := k.listeners[e.Type()]
		if !ok {
			slog.Warn("no listeners for event with matching type", slog.Any("type", e.Type()))
			continue
		}

		for _, r := range runners {
			if r.UpdateValue(e) {
				go func(job runner) {
					err := job.Run(e.Context(ctx), k.DependencyProvider())
					if err != nil {
						slog.Warn("job failed", slog.Any("error", err))
					}
				}(r)
			} else {
				slog.Warn("mismatched event and type, there may be a conflict")
			}
		}

	}
}

func (k *Kernel) Dispatch(ctx context.Context, e event.Event) error {
	return k.queue.Push(e)
}
