package dbpubsub_test

import (
	"testing"

	"abibby.com/salusa/database"
	"abibby.com/salusa/di"
	"abibby.com/salusa/internal/test"
	"abibby.com/salusa/pubsub"
	"abibby.com/salusa/pubsub/dbpubsub"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

type registerDeps struct {
	PS    pubsub.PubSub `inject:""`
	Topic pubsub.Topic  `inject:"default"`
}

func TestRegister(t *testing.T) {
	test.Run(t, "register", func(t *testing.T, tx *sqlx.Tx) {
		if tx.DriverName() == "mysql" {
			t.SkipNow()
		}
		for _, m := range dbpubsub.Migrations {
			err := m.Up.Run(t.Context(), tx)
			if !assert.NoError(t, err) {
				return
			}
		}

		ctx := di.TestDependencyProviderContext()
		di.RegisterSingleton(ctx, func() database.Update {
			return database.Update(func(f func(tx *sqlx.Tx) error) error {
				return f(tx)
			})
		})

		dbpubsub.Register(ctx)

		d, err := di.Resolve[registerDeps](ctx)
		if !assert.NoError(t, err) {
			return
		}
		assert.NotNil(t, d.PS)

		bg := t.Context()
		if !assert.NoError(t, d.Topic.Enqueue(bg, []byte("hello"))) {
			return
		}

		m, err := d.Topic.Dequeue(bg)
		if !assert.NoError(t, err) {
			return
		}
		assert.Equal(t, []byte("hello"), m.Data())

		d2, err := di.Resolve[registerDeps](ctx)
		if !assert.NoError(t, err) {
			return
		}
		assert.Same(t, d.PS, d2.PS)
	})
}
