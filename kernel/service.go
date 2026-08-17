package kernel

import "context"

type Service interface {
	Run(ctx context.Context) error
	Name() string
}

type Restarter interface {
	Restart()
}

type ServiceFunc func() error

func (s ServiceFunc) Run(ctx context.Context) error {
	return s()
}

type ServiceFuncRestart func() error

func (s ServiceFuncRestart) Run(ctx context.Context) error {
	return s()
}
func (s ServiceFuncRestart) Restart() {
}
