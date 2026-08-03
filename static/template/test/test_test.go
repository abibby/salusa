package test

import (
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetupTestDB(t *testing.T) {
	db, err := setupTestDB("sqlite", ":memory:")
	require.NoError(t, err)
	require.NotNil(t, db)
	require.NoError(t, db.Close())

	_, err = setupTestDB("unknown-driver", ":memory:")
	assert.Error(t, err)

	_, err = setupTestDB("sqlite", "/nonexistent-dir/does-not-exist.sqlite")
	assert.Error(t, err)
}

func TestRun(t *testing.T) {
	ran := false
	Run(t, "runner", func(t *testing.T, tx *sqlx.Tx) {
		ran = true
	})
	assert.True(t, ran)
}
