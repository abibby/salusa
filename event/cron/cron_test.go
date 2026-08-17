package cron

import (
	"context"
	"io"
	"log/slog"
	"reflect"
	"testing"
	"time"

	"github.com/abibby/salusa/event"
	"github.com/stretchr/testify/assert"
)

type fakeQueue struct {
	ch chan Event
}

func (q *fakeQueue) Push(ctx context.Context, e event.Event) error {
	q.ch <- e.(Event)
	return nil
}
func (q *fakeQueue) Pop(ctx context.Context, events map[event.EventType]reflect.Type) (event.Event, error) {
	return nil, nil
}

type testEvent struct {
	CronEvent
}

func (e *testEvent) Type() event.EventType {
	return "test"
}

func TestCronEvent_SetTime(t *testing.T) {
	e := &testEvent{}
	now := time.Now()
	e.SetTime(now)
	assert.Equal(t, now, e.Time)
}

func TestServiceAndName(t *testing.T) {
	s := Service()
	assert.Equal(t, "cron-service", s.Name())
	assert.Empty(t, s.events)
}

func TestSchedule(t *testing.T) {
	s := Service()
	e1 := &testEvent{}
	e2 := &testEvent{}

	result := s.Schedule("* * * * *", e1)
	assert.Same(t, s, result)
	assert.Len(t, s.events["* * * * *"], 1)

	s.Schedule("* * * * *", e2)
	assert.Len(t, s.events["* * * * *"], 2)

	s.Schedule("*/5 * * * *", e1)
	assert.Len(t, s.events["*/5 * * * *"], 1)
}

func TestRun(t *testing.T) {
	q := &fakeQueue{ch: make(chan Event, 1)}
	s := Service()
	s.Queue = q
	s.Logger = slog.New(slog.NewTextHandler(io.Discard, nil))

	s.Schedule("* * * * *", &testEvent{})
	s.Schedule("bad-spec", &testEvent{})

	err := s.Run(context.Background())
	assert.NoError(t, err)
}

func TestRunFiresEvent(t *testing.T) {
	q := &fakeQueue{ch: make(chan Event, 1)}
	s := Service()
	s.Queue = q
	s.Logger = slog.New(slog.NewTextHandler(io.Discard, nil))

	s.Schedule("@every 100ms", &testEvent{})

	err := s.Run(context.Background())
	assert.NoError(t, err)

	select {
	case e := <-q.ch:
		te, ok := e.(*testEvent)
		assert.True(t, ok)
		assert.False(t, te.Time.IsZero())
	case <-time.After(2 * time.Second):
		t.Fatal("expected cron event to fire")
	}
}
