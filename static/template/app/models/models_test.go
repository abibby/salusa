package models

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserQuery(t *testing.T) {
	q := UserQuery(context.Background())
	require.NotNil(t, q)
	assert.Equal(t, "users", q.GetTable())
}

func TestFooQuery(t *testing.T) {
	q := FooQuery(context.Background())
	require.NotNil(t, q)
	assert.Equal(t, "foos", q.GetTable())
}
