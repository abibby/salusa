package openapidocdi

import (
	"context"
	"testing"

	"github.com/abibby/salusa/di"
	"github.com/abibby/salusa/kernel"
	"github.com/abibby/salusa/openapidoc"
	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	ctx := di.TestDependencyProviderContext()
	k := kernel.New()
	di.Register(ctx, func(ctx context.Context, tag string) (*kernel.Kernel, error) {
		return k, nil
	})

	err := Register(ctx)
	assert.NoError(t, err)

	api, err := di.Resolve[openapidoc.APIDocer](ctx)
	assert.NoError(t, err)
	assert.Same(t, k, api)
}
