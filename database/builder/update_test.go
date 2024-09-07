package builder_test

import (
	"context"
	"testing"

	"github.com/abibby/salusa/database/builder"
	"github.com/abibby/salusa/internal/test"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestUpdater(t *testing.T) {
	test.QueryTest(t, []test.Case{
		{
			Name:             "Update all",
			Builder:          NewTestBuilder().Updater(builder.Updates{"id": 1}),
			ExpectedSQL:      `UPDATE "foos" SET "id"=?`,
			ExpectedBindings: []any{1},
		},
		{
			Name:             "Update all multi",
			Builder:          NewTestBuilder().Updater(builder.Updates{"id": 1, "foo": "bar"}),
			ExpectedSQL:      `UPDATE "foos" SET "foo"=?, "id"=?`,
			ExpectedBindings: []any{"bar", 1},
		},
		{
			Name:             "Update where",
			Builder:          NewTestBuilder().Where("id", "=", 5).Updater(builder.Updates{"id": 1}),
			ExpectedSQL:      `UPDATE "foos" SET "id"=? WHERE "id" = ?`,
			ExpectedBindings: []any{1, 5},
		},
		{
			Name:             "Update where multi",
			Builder:          NewTestBuilder().Where("id", "=", 5).Updater(builder.Updates{"id": 1, "foo": "bar"}),
			ExpectedSQL:      `UPDATE "foos" SET "foo"=?, "id"=? WHERE "id" = ?`,
			ExpectedBindings: []any{"bar", 1, 5},
		},
	})
}

func TestUpdate(t *testing.T) {
	test.Run(t, "update", func(t *testing.T, tx *sqlx.Tx) {
		const insert = "INSERT INTO foos (id, name) values (?,?)"
		_, err := tx.ExecContext(context.Background(), insert, 1, "test1")
		assert.NoError(t, err)
		_, err = tx.ExecContext(context.Background(), insert, 2, "test2")
		assert.NoError(t, err)

		err = builder.From[*test.Foo]().Where("id", "=", 1).Update(tx, builder.Updates{"name": "new test1"})
		assert.NoError(t, err)

		foos, err := builder.From[*test.Foo]().Get(tx)
		assert.NoError(t, err)
		assert.Len(t, foos, 2)
		assert.Equal(t, 1, foos[0].ID)
		assert.Equal(t, "new test1", foos[0].Name)
	})
}
