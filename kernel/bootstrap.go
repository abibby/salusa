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

	err := k.registerConfig(ctx)
	if err != nil {
		return fmt.Errorf("Kernel.Bootstrap: registerConfig: %w", err)
	}

	for i, b := range k.bootstrap {
		err = b(ctx)
		if err != nil {
			return fmt.Errorf("Kernel.Bootstrap: step %d: %w", i, err)
		}
	}

	return nil
}
