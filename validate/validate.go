package validate

import (
	"context"
	"errors"
)

type Validator interface {
	Validate(ctx context.Context) error
}

func Append(ctx context.Context, err error, v Validator) error {
	return errors.Join(append(expand(err), expand(v.Validate(ctx))...)...)
}
func expand(err error) []error {
	if err == nil {
		return []error{}
	}
	if u, ok := err.(interface{ Unwrap() []error }); ok {
		return u.Unwrap()
	}
	return []error{err}
}
