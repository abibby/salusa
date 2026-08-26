package mixins_test

import (
	"testing"

	"abibby.com/salusa/database/builder"
	"abibby.com/salusa/database/model"
	"abibby.com/salusa/database/model/mixins"
	"abibby.com/salusa/internal/test"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestSoftDeletes(t *testing.T) {
	test.Run(t, "soft_delete", func(t *testing.T, tx *sqlx.Tx) {
		foo := &test.FooSoftDelete{}
		err := model.Save(tx, foo)
		assert.NoError(t, err)

		err = builder.From[*test.FooSoftDelete]().Where("id", "=", 1).Delete(tx)
		assert.NoError(t, err)

		foos, err := builder.From[*test.FooSoftDelete]().Get(tx)
		assert.NoError(t, err)
		assert.Len(t, foos, 0)

		foos, err = builder.From[*test.FooSoftDelete]().WithoutGlobalScope(mixins.SoftDeleteScope).Dump().Get(tx)
		if assert.NoError(t, err) {
			assert.Len(t, foos, 1)
			assert.Equal(t, foo.ID, foos[0].ID)
			assert.NotNil(t, foos[0].DeletedAt)
		}
	})
	//content
}
