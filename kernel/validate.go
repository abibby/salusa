package kernel

import (
	"context"
	"errors"

	"github.com/abibby/salusa/di"
)

func (k *Kernel) Validate(ctx context.Context) error {
	errs := []error{}
	for _, s := range k.services {
		if v, ok := s.(Validator); ok {
			err := v.Validate(ctx)
			if err != nil {
				errs = append(errs, err)
			}
		}

		err := di.Validate(ctx, s)
		if err != nil {
			errs = append(errs, err)
		}
	}

	h := k.rootHandler(ctx)
	if v, ok := h.(Validator); ok {
		err := v.Validate(ctx)
		if err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}
