package kernel

type Service interface {
	Run(k *Kernel) error
	Restart() bool
}

type ServiceFunc func(k *Kernel) error

func (s ServiceFunc) Run(k *Kernel) error {
	return s(k)
}
func (s ServiceFunc) Restart() bool {
	return false
}

type ServiceFuncRestart func(k *Kernel) error

func (s ServiceFuncRestart) Run(k *Kernel) error {
	return s(k)
}
func (s ServiceFuncRestart) Restart() bool {
	return true
}
