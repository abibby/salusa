package event

import (
	"context"
	"log/slog"
	"reflect"

	"github.com/abibby/salusa/di"
	"github.com/abibby/salusa/kernel"
)

type JobHandler[E Event] interface {
	Handle(ctx context.Context, event E) error
}

type runner interface {
	UpdateValue(v Event) bool
	Run(ctx context.Context, dp *di.DependencyProvider) error
	EventType() reflect.Type
}
type Listener struct {
	eventType EventType
	runner    runner
}

type job[E Event] struct {
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

func (j *job[E]) UpdateValue(v Event) bool {
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

func NewListener[H JobHandler[E], E Event]() *Listener {
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

type EventService struct {
	Queue  Queue                  `inject:""`
	Logger *slog.Logger           `inject:""`
	DP     *di.DependencyProvider `inject:""`

	listeners map[EventType][]runner
}

var _ kernel.Service = (*EventService)(nil)

func Service() *EventService {
	return &EventService{
		listeners: map[EventType][]runner{},
	}
}

func (s *EventService) Name() string {
	return "event-service"
}

func (s *EventService) Run(ctx context.Context, k *kernel.Kernel) error {

	events := map[EventType]reflect.Type{}
	for eventType, runners := range s.listeners {
		events[eventType] = runners[0].EventType()
	}

	for {
		e, err := s.Queue.Pop(events)
		if err != nil {
			s.Logger.Warn("could not pop event off queue", slog.Any("error", err))
			continue
		}
		runners, ok := s.listeners[e.Type()]
		if !ok {
			s.Logger.Warn("no listeners for event with matching type", slog.Any("type", e.Type()))
			continue
		}

		for _, r := range runners {
			if r.UpdateValue(e) {
				go func(job runner) {
					err := job.Run(e.Context(ctx), s.DP)
					if err != nil {
						s.Logger.Warn("job failed", slog.Any("error", err))
					}
				}(r)
			} else {
				s.Logger.Warn("mismatched event and type, there may be a conflict")
			}
		}

	}
}

func (s *EventService) Add(l *Listener) *EventService {
	jobs, ok := s.listeners[l.eventType]
	if !ok {
		jobs = []runner{}
	}
	s.listeners[l.eventType] = append(jobs, l.runner)
	return s
}
