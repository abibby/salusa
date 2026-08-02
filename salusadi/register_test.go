package salusadi_test

import (
	"context"
	"testing"

	"github.com/abibby/salusa/auth"
	"github.com/abibby/salusa/database/migrate"
	"github.com/abibby/salusa/di"
	"github.com/abibby/salusa/salusadi"
	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	migrations := migrate.New()
	register := salusadi.Register[*auth.UsernameUser](migrations)

	ctx := di.TestDependencyProviderContext()
	err := register(ctx)
	assert.NoError(t, err)
}

func TestRegisterWithNilMigrations(t *testing.T) {
	register := salusadi.Register[*auth.UsernameUser](nil)

	ctx := di.TestDependencyProviderContext()
	err := register(ctx)
	assert.NoError(t, err)
}

var _ = context.Background
