package builder_test

import (
	"context"
	"testing"

	"github.com/abibby/salusa/database/builder"
	"github.com/abibby/salusa/internal/test"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestDeleter(t *testing.T) {
	test.QueryTest(t, []test.Case{
		{
			Name:             "delete all",
			Builder:          NewTestBuilder().Deleter(),
			ExpectedSQL:      "DELETE FROM \"foos\"",
			ExpectedBindings: []any{},
		},
		{
			Name:             "delete where",
			Builder:          NewTestBuilder().Where("id", "=", 5).Deleter(),
			ExpectedSQL:      "DELETE FROM \"foos\" WHERE \"id\" = ?",
			ExpectedBindings: []any{5},
		},
	})
}

func TestDelete(t *testing.T) {
	test.Run(t, "delete", func(t *testing.T, tx *sqlx.Tx) {
		const insert = "INSERT INTO foos (id, name) values (?,?)"
		_, err := tx.ExecContext(context.Background(), insert, 1, "test1")
		assert.NoError(t, err)
		_, err = tx.ExecContext(context.Background(), insert, 2, "test2")
		assert.NoError(t, err)

		err = builder.From[*test.Foo]().Where("id", "=", 1).Delete(tx)
		assert.NoError(t, err)

		foos, err := builder.From[*test.Foo]().Get(tx)
		assert.NoError(t, err)
		assert.Len(t, foos, 1)
		assert.Equal(t, 2, foos[0].ID)
	})
}
