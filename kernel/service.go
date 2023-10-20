package kernel

import "context"

type Service interface {
	Run(ctx context.Context, k *Kernel) error
	Name() string
}

type Restarter interface {
	Restart()
}

type ServiceFunc func(k *Kernel) error

func (s ServiceFunc) Run(ctx context.Context, k *Kernel) error {
	return s(k)
}
func (s ServiceFunc) Restart() bool {
	return false
}

type ServiceFuncRestart func(k *Kernel) error

func (s ServiceFuncRestart) Run(ctx context.Context, k *Kernel) error {
	return s(k)
}
