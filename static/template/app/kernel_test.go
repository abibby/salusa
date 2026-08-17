package app

import (
	"testing"

	"github.com/abibby/salusa/database/dialects/sqlite"
	"github.com/abibby/salusa/di"
	"github.com/abibby/salusa/email/emailtest"
	"github.com/abibby/salusa/kernel"
	"github.com/abibby/salusa/static/template/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestKernel(t *testing.T) {
	assert.NotNil(t, Kernel)
}

func TestKernelBootstrap(t *testing.T) {
	ctx := di.TestDependencyProviderContext()

	cfg := &config.Config{
		Port:     443,
		BasePath: "https://example.test",

		Database: sqlite.NewConfig(":memory:"),
		Mail:     emailtest.NewTestMailerConfig(),
	}

	k := kernel.Config(func() *config.Config { return cfg })(Kernel)

	err := k.Bootstrap(ctx)
	require.NoError(t, err)
}
