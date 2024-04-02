package kernel

import (
	"context"
	"errors"
	"fmt"
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
	for i, b := range k.bootstrap {
		err = b(ctx)
		if err != nil {
			return fmt.Errorf("error %d: %w", i, err)
		}
	}
	for _, b := range k.postBootstrap {
		b()
	}
	return nil
}
