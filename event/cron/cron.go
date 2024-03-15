package cron

import (
	"context"
	"log/slog"
	"time"

	"github.com/abibby/salusa/event"
	"github.com/abibby/salusa/kernel"
	"github.com/robfig/cron/v3"
)

type Event interface {
	event.Event
	SetTime(t time.Time)
}

type CronEvent struct {
	Time time.Time
}

func (b *CronEvent) SetTime(t time.Time) {
	b.Time = t
}

type CronService struct {
	events map[string][]Event
}

func Service() *CronService {
	return &CronService{
		events: map[string][]Event{},
	}
}

func (c *CronService) Run(ctx context.Context, k *kernel.Kernel) error {
	runner := cron.New()
	for spec, events := range c.events {
		for _, e := range events {
			_, err := runner.AddFunc(spec, func() {
				e.SetTime(time.Now())
				err := k.Dispatch(ctx, e)
				if err != nil {
					k.Logger(ctx).Error("failed to dispatch event", slog.Any("error", err))
				}
			})
			if err != nil {
				k.Logger(ctx).Error("failed to start cron listener", slog.Any("error", err))
			}
		}
	}
	runner.Start()

	return nil
}
func (c *CronService) Name() string {
	return "cron-service"
}

func (c *CronService) Schedule(cron string, e Event) *CronService {
	jobs, ok := c.events[cron]
	if !ok {
		jobs = []Event{}
	}
	c.events[cron] = append(jobs, e)
	return c
}
