package test_test

import (
	"testing"

	"github.com/abibby/salusa/static/template/test"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRun(t *testing.T) {
	test.Run(t, "runs with a transaction", func(t *testing.T, tx *sqlx.Tx) {
		assert.NotNil(t, tx)
	})
}

func TestHelpers(t *testing.T) {
	require.NotNil(t, test.Run)
	require.NotNil(t, test.RunBenchmark)
	require.NotNil(t, test.Kernel)
}
