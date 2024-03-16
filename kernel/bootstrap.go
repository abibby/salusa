package kernel

import (
	"context"
	"errors"
)

var (
	ErrAlreadyBootstrapped = errors.New("kernel already bootstrapped")
)

func (k *Kernel) Bootstrap(ctx context.Context) error {
	if k.bootstrapped {
		return ErrAlreadyBootstrapped
	}
	k.bootstrapped = true

	var err error
	for _, p := range k.providers {
		p(k.dp)
	}
	for _, b := range k.bootstrap {
		err = b(ctx, k)
		if err != nil {
			return err
		}
	}
	for _, b := range k.postBootstrap {
		b()
	}
	return nil
}
