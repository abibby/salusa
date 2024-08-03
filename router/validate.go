package router

import (
	"context"

	"github.com/abibby/salusa/validate"
)

func (r *Router) Validate(ctx context.Context) error {
	var err error
	for _, route := range r.Routes() {
		if v, ok := route.handler.(validate.Validator); ok {
			err = validate.Append(ctx, err, v)
		}
	}
	return err
}
