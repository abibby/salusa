package event

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"reflect"
	"testing"
	"time"

	"github.com/abibby/salusa/di"
	"github.com/abibby/salusa/salusaconfig"
	"github.com/stretchr/testify/assert"
)

func TestChannelQueue(t *testing.T) {
	t.Run("push pop", func(t *testing.T) {
		q := NewChannelQueue()
		ev := &TestEvent1{Foo: "bar"}
		err := q.Push(ev)
		assert.NoError(t, err)
		e, err := q.Pop(map[EventType]reflect.Type{
			ev.Type(): reflect.TypeOf(&TestEvent1{}),
		})
		assert.NoError(t, err)
		assert.Equal(t, ev, e)
	})

	t.Run("config", func(t *testing.T) {
		c := NewChannelQueueConfig()
		assert.IsType(t, &ChannelQueue{}, c.Queue())
	})

	t.Run("register", func(t *testing.T) {
		ctx := di.TestDependencyProviderContext()
		err := RegisterChannelQueue(ctx)
		assert.NoError(t, err)
		q, err := di.Resolve[Queue](ctx)
		assert.NoError(t, err)
		assert.NotNil(t, q)
	})
}

type queueTestConfig struct{}

func (queueTestConfig) GetHTTPPort() int {
	return 8080
}
func (queueTestConfig) GetBaseURL() string {
	return "https://example.com"
}
func (queueTestConfig) QueueConfig() Config {
	return NewChannelQueueConfig()
}

type plainConfig struct{}

func (plainConfig) GetHTTPPort() int {
	return 8080
}
func (plainConfig) GetBaseURL() string {
	return "https://example.com"
}

func TestQueueRegister(t *testing.T) {
	t.Run("with queue config", func(t *testing.T) {
		ctx := di.TestDependencyProviderContext()
		di.RegisterSingleton(ctx, func() salusaconfig.Config {
			return queueTestConfig{}
		})

		err := Register(ctx)
		assert.NoError(t, err)
		q, err := di.Resolve[Queue](ctx)
		assert.NoError(t, err)
		assert.NotNil(t, q)
	})

	t.Run("without queue config", func(t *testing.T) {
		ctx := di.TestDependencyProviderContext()
		di.RegisterSingleton(ctx, func() salusaconfig.Config {
			return plainConfig{}
		})

		err := Register(ctx)
		assert.NoError(t, err)
		_, err = di.Resolve[Queue](ctx)
		assert.Error(t, err)
	})
}

var testHandlerDone chan string

type testEventHandler struct{}

func (testEventHandler) Handle(ctx context.Context, e *TestEvent1) error {
	testHandlerDone <- e.Foo
	return nil
}

type testErrorHandler struct{}

func (testErrorHandler) Handle(ctx context.Context, e *TestEvent1) error {
	return fmt.Errorf("handler boom")
}

type testDep struct{ V int }

var fillableResult string

type fillableHandler struct {
	Dep *testDep `inject:""`
}

func (h *fillableHandler) Handle(ctx context.Context, e *TestEvent1) error {
	fillableResult = fmt.Sprintf("%s:%d", e.Foo, h.Dep.V)
	return nil
}

func TestHandler(t *testing.T) {
	t.Run("run", func(t *testing.T) {
		testHandlerDone = make(chan string, 1)
		h := &handler[*TestEvent1]{
			value:       &TestEvent1{Foo: "bar"},
			handlerType: reflect.TypeFor[testEventHandler](),
		}
		err := h.Run(context.Background(), di.NewDependencyProvider())
		assert.NoError(t, err)
		assert.Equal(t, "bar", <-testHandlerDone)
	})

	t.Run("run fillable", func(t *testing.T) {
		dp := di.NewDependencyProvider()
		dp.Register(di.NewSingletonFactory(&testDep{V: 7}))
		h := &handler[*TestEvent1]{
			value:       &TestEvent1{Foo: "bar"},
			handlerType: reflect.TypeFor[*fillableHandler](),
		}
		err := h.Run(context.Background(), dp)
		assert.NoError(t, err)
		assert.Equal(t, "bar:7", fillableResult)
	})

	t.Run("run fillable error", func(t *testing.T) {
		h := &handler[*TestEvent1]{
			value:       &TestEvent1{Foo: "bar"},
			handlerType: reflect.TypeFor[*fillableHandler](),
		}
		err := h.Run(context.Background(), di.NewDependencyProvider())
		assert.ErrorIs(t, err, di.ErrNotRegistered)
	})

	t.Run("update value", func(t *testing.T) {
		h := &handler[*TestEvent1]{}
		assert.False(t, h.UpdateValue(&TestEvent2{}))
		assert.True(t, h.UpdateValue(&TestEvent1{Foo: "x"}))
		assert.Equal(t, "x", h.value.Foo)
	})

	t.Run("event type", func(t *testing.T) {
		h := &handler[*TestEvent1]{}
		assert.Equal(t, reflect.TypeOf(&TestEvent1{}), h.EventType())
	})
}

func TestService(t *testing.T) {
	s := Service(NewListener[*testEventHandler, *TestEvent1]())
	assert.Equal(t, "event-service", s.Name())
	assert.Len(t, s.listeners, 1)
	assert.Len(t, s.listeners[(&TestEvent1{}).Type()], 1)
}

func TestNewListener(t *testing.T) {
	l := NewListener[*testEventHandler, *TestEvent1]()
	assert.Equal(t, EventType("test-event:1"), l.eventType)
	assert.NotNil(t, l.runner)
}

type queueItem struct {
	err   error
	event Event
}

type testQueue struct {
	items []queueItem
	idx   int
	block chan struct{}
}

func (q *testQueue) Push(e Event) error {
	return nil
}

func (q *testQueue) Pop(events map[EventType]reflect.Type) (Event, error) {
	if q.idx < len(q.items) {
		item := q.items[q.idx]
		q.idx++
		return item.event, item.err
	}
	<-q.block
	return nil, nil
}

func TestEventServiceRun(t *testing.T) {
	testHandlerDone = make(chan string, 1)
	mismatchListener := &Listener{
		eventType: (&TestEvent1{}).Type(),
		runner:    &handler[*TestEvent2]{},
	}
	s := Service(
		NewListener[*testEventHandler, *TestEvent1](),
		NewListener[*testErrorHandler, *TestEvent1](),
		mismatchListener,
	)
	s.Queue = &testQueue{
		items: []queueItem{
			{err: fmt.Errorf("pop boom")},
			{event: &TestEvent2{}},
			{event: &TestEvent1{Foo: "ok"}},
		},
		block: make(chan struct{}),
	}
	s.Logger = slog.New(slog.NewTextHandler(io.Discard, nil))
	s.DP = di.NewDependencyProvider()

	go func() {
		_ = s.Run(context.Background())
	}()

	select {
	case v := <-testHandlerDone:
		assert.Equal(t, "ok", v)
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for event to be handled")
	}
}
