package kernel

import (
	"context"
	"reflect"

	"log/slog"

	"github.com/abibby/salusa/event"
)

type runner interface {
	UpdateValue(v event.Event) bool
	Run(ctx context.Context) error
}
type Listener struct {
	eventType event.EventType
	runner    runner
}

type job[E event.Event] struct {
	value    E
	callback func(ctx context.Context, event E) error
}

func (j *job[E]) Run(ctx context.Context) error {
	return j.callback(ctx, j.value)
}
func (j *job[E]) UpdateValue(v event.Event) bool {
	ev, ok := v.(E)
	if !ok {
		return false
	}
	j.value = ev
	return true
}

func NewListener[E event.Event](cb func(ctx context.Context, event E) error) *Listener {
	var e E
	return &Listener{
		eventType: e.Type(),
		runner: &job[E]{
			value:    e,
			callback: cb,
		},
	}
}

func (k *Kernel) RunListeners(ctx context.Context) {
	events := map[event.EventType]reflect.Type{}
	for eventType, runners := range k.listeners {
		f, ok := reflect.TypeOf(runners[0]).Elem().FieldByName("callback")
		if ok {
			events[eventType] = f.Type.In(1)
		}
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
					err := job.Run(e.Context(ctx))
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
