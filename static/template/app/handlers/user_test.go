package handlers_test

import (
	"context"
	"testing"

	"github.com/abibby/salusa/kernel"
	"github.com/abibby/salusa/static/template/app"
	"github.com/abibby/salusa/static/template/app/handlers"
	"github.com/abibby/salusa/static/template/app/models/factories"
	"github.com/abibby/salusa/static/template/test"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestUser(t *testing.T) {
	app.Kernel = kernel.New()
	test.Run(t, "create user", func(t *testing.T, tx *sqlx.Tx) {
		resp, err := handlers.UserCreate.Test(&handlers.CreateUserRequest{
			Username: "name",
			Password: []byte("word"),
			Tx:       tx,
			Ctx:      context.Background(),
		})

		assert.NoError(t, err)
		assert.Equal(t, "name", resp.User.Username)
		assert.NotZero(t, resp.User.ID)

	})

	test.Run(t, "get user", func(t *testing.T, tx *sqlx.Tx) {
		u := factories.UserFactory.Create(tx)

		resp, err := handlers.UserGet.Test(&handlers.GetUserRequest{
			User: u,
		})

		assert.NoError(t, err)
		assert.Equal(t, u, resp.User)

	})
}
