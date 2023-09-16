package kernel

import (
	"log"

	"github.com/abibby/salusa/event"
)

type runner interface {
	UpdateValue(v event.Event) bool
	Run() error
}
type Listener struct {
	eventType event.EventType
	runner    runner
}

type job[E event.Event] struct {
	value    E
	callback func(event E) error
}

func (j *job[E]) Run() error {
	return j.callback(j.value)
}
func (j *job[E]) UpdateValue(v event.Event) bool {
	ev, ok := v.(E)
	if !ok {
		return false
	}
	j.value = ev
	return true
}

func NewListener[E event.Event](cb func(event E) error) *Listener {
	var e E
	return &Listener{
		eventType: e.Type(),
		runner: &job[E]{
			value:    e,
			callback: cb,
		},
	}
}

func (k *Kernel) runListeners() {
	for {
		e := k.queue.Pop()
		runners, ok := k.listeners[e.Type()]
		if !ok {
			log.Printf("no listeners for event with type %s", e.Type())
			continue
		}

		for _, r := range runners {
			if r.UpdateValue(e) {
				go func(l runner) {
					err := l.Run()
					if err != nil {
						log.Print(err)
					}
				}(r)
			} else {
				log.Printf("mismatched event and type, there may be a conflict")
			}
		}

	}
}

func (k *Kernel) Dispatch(e event.Event) {
	k.queue.Push(e)
}

func Dispatch(e event.Event) {
	defaultKernel.Dispatch(e)
}
