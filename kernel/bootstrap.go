package kernel

import "context"

func (k *Kernel) Bootstrap(ctx context.Context) error {
	var err error
	for _, b := range k.bootstrap {
		err = b(ctx)
		if err != nil {
			return err
		}
	}
	for _, b := range k.postBootstrap {
		b()
	}
	return nil
}
