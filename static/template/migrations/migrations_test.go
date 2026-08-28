package migrations_test

import (
	"testing"

	"abibby.com/salusa/static/template/migrations"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUse(t *testing.T) {
	m := migrations.Use()
	require.NotNil(t, m)
	assert.NotNil(t, m)
}
