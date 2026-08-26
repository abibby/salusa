package dbpubsub_test

import (
	"testing"

	"abibby.com/salusa/database"
	"abibby.com/salusa/internal/test"
	"abibby.com/salusa/pubsub"
	"abibby.com/salusa/pubsub/dbpubsub"
	"abibby.com/salusa/pubsub/pubsubtest"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestStandard(t *testing.T) {
	pubsubtest.RunStandardTests(t, func(t *testing.T, name string, fn func(*testing.T, pubsub.PubSub)) {
		test.Run(t, name, func(t *testing.T, tx *sqlx.Tx) {
			if tx.DriverName() == "mysql" {
				t.SkipNow()
			}
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
