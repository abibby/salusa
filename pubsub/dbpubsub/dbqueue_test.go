package dbpubsub_test

import (
	"testing"

	"github.com/abibby/salusa/database"
	"github.com/abibby/salusa/internal/test"
	"github.com/abibby/salusa/pubsub"
	"github.com/abibby/salusa/pubsub/dbpubsub"
	"github.com/abibby/salusa/pubsub/pubsubtest"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestStandard(t *testing.T) {
	pubsubtest.RunStandardTests(t, func(t *testing.T, name string, fn func(*testing.T, pubsub.PubSub)) {
		test.Run(t, name, func(t *testing.T, tx *sqlx.Tx) {
			for _, m := range dbpubsub.Migrations {
				err := m.Up.Run(t.Context(), tx)
				if !assert.NoError(t, err) {
					return
				}
			}
			p := dbpubsub.New(database.Update(func(f func(tx *sqlx.Tx) error) error {
				return f(tx)
			}))
			fn(t, p)
		})
	})
}
