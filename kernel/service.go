package kernel

type Service interface {
	Run(k *Kernel) error
	Restart() bool
}
