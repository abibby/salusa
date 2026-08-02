package providers

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegister(t *testing.T) {
	called := false
	Add(func(ctx context.Context) {
		called = true
	})

	err := Register(context.Background())
	require.NoError(t, err)
	assert.True(t, called)
}
