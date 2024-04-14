package router

import (
	"context"
	"errors"

	"github.com/abibby/salusa/kernel"
)

func (r *Router) Validate(ctx context.Context) error {
	errs := []error{}
	for _, route := range r.Routes() {
		if v, ok := route.handler.(kernel.Validator); ok {
			err := v.Validate(ctx)
			if err != nil {
				errs = append(errs, err)
			}
		}
	}
	return errors.Join(errs...)
}
