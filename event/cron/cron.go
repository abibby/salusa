package cron

import (
	"time"

	"github.com/abibby/salusa/event"
	"github.com/abibby/salusa/kernel"
	"github.com/robfig/cron/v3"
)

type Event interface {
	event.Event
	SetTime(t time.Time)
}

type BaseEvent struct {
	Time time.Time
}

func (b *BaseEvent) SetTime(t time.Time) {
	b.Time = t
}

type CronService struct {
	events map[string][]Event
}

func NewService() *CronService {
	return &CronService{
		events: map[string][]Event{},
	}
}

func (c *CronService) Run(k *kernel.Kernel) error {
	runner := cron.New()
	for spec, events := range c.events {
		for _, e := range events {
			runner.AddFunc(spec, func() {
				e.SetTime(time.Now())
				k.Dispatch(e)
			})
		}
	}
	runner.Start()

	return nil
}
func (c *CronService) Restart() bool {
	return false
}

func (c *CronService) Schedule(cron string, e Event) *CronService {
	jobs, ok := c.events[cron]
	if !ok {
		jobs = []Event{}
	}
	c.events[cron] = append(jobs, e)
	return c
}
