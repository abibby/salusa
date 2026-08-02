package kernel

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/abibby/salusa/di"
	"github.com/abibby/salusa/router"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
)

type runTestConfig struct {
	port int
}

func (c *runTestConfig) GetHTTPPort() int   { return c.port }
func (c *runTestConfig) GetBaseURL() string { return "http://localhost" }

func newRunTestKernel(t *testing.T, port int) (*Kernel, context.Context) {
	t.Helper()
	k := New(
		Config(func() *runTestConfig { return &runTestConfig{port: port} }),
		InitRoutes(func(r *router.Router) {
			r.Get("/ping", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte("pong"))
			}))
		}),
	)
	ctx := di.TestDependencyProviderContext()
	err := k.Bootstrap(ctx)
	assert.NoError(t, err)
	return k, ctx
}

func TestHandlerWithMiddleware(t *testing.T) {
	k, _ := newRunTestKernel(t, 0)
	h := k.handlerWithMiddleware()
	assert.NotNil(t, h)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/ping", http.NoBody)
	h.ServeHTTP(w, r)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestKernelHttpServer(t *testing.T) {
	k, ctx := newRunTestKernel(t, 8123)
	server := k.HttpServer(ctx)
	assert.Equal(t, ":8123", server.Addr)
	assert.NotNil(t, server.Handler)

	resolved, err := di.Resolve[*http.Server](ctx)
	assert.NoError(t, err)
	assert.Same(t, server, resolved)
}

func TestRunHttpServer(t *testing.T) {
	k, ctx := newRunTestKernel(t, 0)

	errCh := make(chan error, 1)
	go func() {
		errCh <- k.RunHttpServer(ctx)
	}()

	var server *http.Server
	assert.Eventually(t, func() bool {
		s, err := di.Resolve[*http.Server](ctx)
		if err != nil {
			return false
		}
		server = s
		return true
	}, time.Second, time.Millisecond*10)

	err := server.Close()
	assert.NoError(t, err)

	select {
	case runErr := <-errCh:
		assert.ErrorIs(t, runErr, http.ErrServerClosed)
	case <-time.After(time.Second):
		t.Fatal("RunHttpServer did not return")
	}
}

type runServicesTestService struct {
	name string
	run  func(ctx context.Context) error
}

func (s *runServicesTestService) Name() string                  { return s.name }
func (s *runServicesTestService) Run(ctx context.Context) error { return s.run(ctx) }

func TestRunServices(t *testing.T) {
	ran := make(chan struct{}, 1)
	svc := &runServicesTestService{
		name: "test-service",
		run: func(ctx context.Context) error {
			ran <- struct{}{}
			return nil
		},
	}
	k := New(Services(svc))
	ctx := di.TestDependencyProviderContext()

	k.RunServices(ctx)

	select {
	case <-ran:
	case <-time.After(time.Second):
		t.Fatal("service did not run")
	}
}

type restartableService struct {
	calls int
	done  chan struct{}
}

func (s *restartableService) Name() string { return "restartable" }
func (s *restartableService) Run(ctx context.Context) error {
	s.calls++
	if s.calls == 1 {
		return fmt.Errorf("fail once")
	}
	close(s.done)
	return nil
}
func (s *restartableService) Restart() {}

func TestRunServicesRestart(t *testing.T) {
	svc := &restartableService{done: make(chan struct{})}
	k := New(Services(svc))
	ctx := di.TestDependencyProviderContext()

	k.RunServices(ctx)

	select {
	case <-svc.done:
	case <-time.After(time.Second):
		t.Fatal("service was not restarted")
	}
	assert.Equal(t, 2, svc.calls)
}

type fillableDep struct{}

type fillableService struct {
	Dep *fillableDep `inject:""`
	ran chan *fillableDep
}

func (s *fillableService) Name() string { return "fillable" }
func (s *fillableService) Run(ctx context.Context) error {
	s.ran <- s.Dep
	return nil
}

func TestRunServicesFillable(t *testing.T) {
	ctx := di.TestDependencyProviderContext()
	di.RegisterSingleton(ctx, func() *fillableDep { return &fillableDep{} })

	svc := &fillableService{ran: make(chan *fillableDep, 1)}
	k := New(Services(svc))

	k.RunServices(ctx)

	select {
	case dep := <-svc.ran:
		assert.NotNil(t, dep)
	case <-time.After(time.Second):
		t.Fatal("service did not run")
	}
}

func TestRunServicesFillableMissingDependency(t *testing.T) {
	ctx := di.TestDependencyProviderContext()

	svc := &fillableService{ran: make(chan *fillableDep, 1)}
	k := New(Services(svc))

	k.RunServices(ctx)

	select {
	case <-svc.ran:
		t.Fatal("service should not have run without its dependency")
	case <-time.After(time.Millisecond * 100):
	}
}

func TestSingles(t *testing.T) {
	k, ctx := newRunTestKernel(t, 0)
	assert.NotPanics(t, func() {
		k.singles(ctx)
	})
}

type closerSingleton struct {
	err error
}

func (c *closerSingleton) Close() error { return c.err }

func TestClose(t *testing.T) {
	k, ctx := newRunTestKernel(t, 0)
	di.RegisterSingleton(ctx, func() *closerSingleton { return &closerSingleton{} })
	_, err := di.Resolve[*closerSingleton](ctx)
	assert.NoError(t, err)

	err = k.Close()
	assert.NoError(t, err)
}

func TestCloseWithError(t *testing.T) {
	k, ctx := newRunTestKernel(t, 0)
	di.RegisterSingleton(ctx, func() *closerSingleton { return &closerSingleton{err: fmt.Errorf("boom")} })
	_, err := di.Resolve[*closerSingleton](ctx)
	assert.NoError(t, err)

	err = k.Close()
	assert.Error(t, err)
}

func TestCloseAndLog(t *testing.T) {
	k, ctx := newRunTestKernel(t, 0)
	di.RegisterSingleton(ctx, func() *closerSingleton { return &closerSingleton{err: fmt.Errorf("boom")} })
	_, err := di.Resolve[*closerSingleton](ctx)
	assert.NoError(t, err)

	assert.NotPanics(t, func() {
		k.closeAndLog(ctx)
	})
}

func TestResourceError(t *testing.T) {
	inner := fmt.Errorf("inner")
	re := &resourceError{err: inner, resource: "thing"}
	assert.Equal(t, "resource thing: inner", re.Error())
	assert.ErrorIs(t, re, inner)
}

func TestRunFetch(t *testing.T) {
	oldArgs := os.Args
	oldCommandLine := pflag.CommandLine
	defer func() {
		os.Args = oldArgs
		pflag.CommandLine = oldCommandLine
	}()
	pflag.CommandLine = pflag.NewFlagSet(oldArgs[0], pflag.ContinueOnError)
	os.Args = []string{oldArgs[0], "--fetch", "/ping"}

	k, ctx := newRunTestKernel(t, 0)

	err := k.Run(ctx)
	assert.NoError(t, err)
}
