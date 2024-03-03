package migrate_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/abibby/salusa/database"
	"github.com/abibby/salusa/database/migrate"
	"github.com/abibby/salusa/database/schema"
	"github.com/abibby/salusa/internal/test"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestMigrations(t *testing.T) {
	test.Run(t, "dont rerun migraitons", func(t *testing.T, tx *sqlx.Tx) {
		m := migrate.New()
		m.Add(&migrate.Migration{
			Name: "1",
			Up: schema.Create("foo", func(b *schema.Blueprint) {
				b.Int("id").Primary()
			}),
		})

		err := m.Up(context.Background(), tx)
		assert.NoError(t, err)
		m.Add(&migrate.Migration{
			Name: "2",
			Up: schema.Table("foo", func(b *schema.Blueprint) {
				b.String("name")
			}),
		})

		err = m.Up(context.Background(), tx)
		assert.NoError(t, err)
	})
	test.Run(t, "failed migrations", func(t *testing.T, tx *sqlx.Tx) {
		m := migrate.New()
		m1 := &migrate.Migration{
			Name: "1",
			Up: schema.Run(func(ctx context.Context, tx database.DB) error {
				return fmt.Errorf("error")
			}),
		}
		m.Add(m1)

		err := m.Up(context.Background(), tx)
		assert.Error(t, err)
		m1.Up = schema.Create("foo", func(b *schema.Blueprint) {
			b.Int("id").Primary()
		})

		m.Add(&migrate.Migration{
			Name: "2",
			Up: schema.Table("foo", func(b *schema.Blueprint) {
				b.String("name")
			}),
		})

		err = m.Up(context.Background(), tx)
		assert.NoError(t, err)
	})
}
