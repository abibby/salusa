package handlers_test

import (
	"context"
	"testing"

	"abibby.com/salusa/database"
	"abibby.com/salusa/static/template/app/handlers"
	"abibby.com/salusa/static/template/app/models"
	"abibby.com/salusa/static/template/test"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserGet(t *testing.T) {
	user := &models.User{}
	resp, err := handlers.UserGet.Run(&handlers.GetUserRequest{
		User: user,
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Same(t, user, resp.User)
}

func TestUserList(t *testing.T) {
	test.Run(t, "user list", func(t *testing.T, tx *sqlx.Tx) {
		resp, err := handlers.UserList.Run(&handlers.ListUserRequest{
			Read: database.Read(func(cb func(tx *sqlx.Tx) error) error {
				return cb(tx)
			}),
			Ctx: context.Background(),
		})
		require.NoError(t, err)
		require.NotNil(t, resp)
		assert.Empty(t, resp.Users)
	})
}
