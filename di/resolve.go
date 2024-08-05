package di

import (
	"context"
)

func Resolve[T any](ctx context.Context) (T, error) {
	dp := GetDependencyProvider(ctx)
	var result T
	err := dp.Fill(ctx, &result)
	if err != nil {
		var zero T
		return zero, err
	}
	return result, err
}
