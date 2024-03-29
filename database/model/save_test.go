package model_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/abibby/salusa/database"
	"github.com/abibby/salusa/database/builder"
	"github.com/abibby/salusa/database/hooks"
	"github.com/abibby/salusa/database/model"
	"github.com/abibby/salusa/internal/test"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestSave_create(t *testing.T) {
	test.Run(t, "create", func(t *testing.T, tx *sqlx.Tx) {
		f := &test.Foo{
			ID:   1,
			Name: "test",
		}
		err := model.Save(tx, f)
		assert.NoError(t, err)

		rows, err := tx.QueryContext(context.Background(), "select id, name from foos")
		assert.NoError(t, err)

		assert.True(t, rows.Next())
		id := 0
		name := ""
		err = rows.Scan(&id, &name)
		assert.NoError(t, err)

		assert.Equal(t, f.ID, id)
		assert.Equal(t, f.Name, name)

		assert.False(t, rows.Next())
	})

	test.Run(t, "autoincrement", func(t *testing.T, tx *sqlx.Tx) {
		f1 := &test.Foo{
			Name: "test",
		}
		err := model.Save(tx, f1)
		assert.NoError(t, err)
		f2 := &test.Foo{
			Name: "test",
		}
		err = model.Save(tx, f2)
		assert.NoError(t, err)

		rows, err := tx.QueryContext(context.Background(), "select id, name from foos")
		assert.NoError(t, err)

		id := 0
		name := ""
		assert.True(t, rows.Next())
		err = rows.Scan(&id, &name)
		assert.NoError(t, err)

		assert.Equal(t, f1.ID, id)
		assert.Equal(t, f1.Name, name)

		assert.True(t, rows.Next())
		err = rows.Scan(&id, &name)
		assert.NoError(t, err)

		assert.Equal(t, f2.ID, id)
		assert.Equal(t, f2.Name, name)

		assert.False(t, rows.Next())
	})
}

func TestSave_update(t *testing.T) {
	test.Run(t, "update", func(t *testing.T, tx *sqlx.Tx) {
		f := &test.Foo{
			ID:   1,
			Name: "test",
		}
		err := model.Save(tx, f)
		assert.NoError(t, err)

		f.Name = "new name"
		err = model.Save(tx, f)
		assert.NoError(t, err)

		rows, err := tx.QueryContext(context.Background(), "select id, name from foos")
		assert.NoError(t, err)

		assert.True(t, rows.Next())
		id := 0
		name := ""
		err = rows.Scan(&id, &name)
		assert.NoError(t, err)

		assert.Equal(t, f.ID, id)
		assert.Equal(t, f.Name, name)

		assert.False(t, rows.Next())
	})
}

func TestSave_model_is_in_database_after_saving(t *testing.T) {
	test.Run(t, "model in database after saving", func(t *testing.T, tx *sqlx.Tx) {
		f := &test.Foo{
			ID: 1,
		}
		err := model.Save(tx, f)
		assert.NoError(t, err)

		assert.True(t, f.InDatabase())
	})
}

func TestSave_autoincrement(t *testing.T) {
	test.Run(t, "autoincrement", func(t *testing.T, tx *sqlx.Tx) {
		f := &test.Foo{}
		err := model.Save(tx, f)
		assert.NoError(t, err)

		assert.Equal(t, f.ID, 1)
	})
	test.Run(t, "autoincrement set id", func(t *testing.T, tx *sqlx.Tx) {
		f := &test.Foo{
			ID: 100,
		}
		err := model.Save(tx, f)
		assert.NoError(t, err)

		assert.Equal(t, f.ID, 100)
	})
}

type FooSaveHookTest struct {
	test.Foo
	saved bool
}

type FooSaveHookTestWrapper struct {
	FooSaveHookTest
}

var _ hooks.AfterSaver = &FooSaveHookTest{}

func (f *FooSaveHookTest) AfterSave(context.Context, database.DB) error {
	f.saved = true
	return nil
}
func (f *FooSaveHookTest) Table() string {
	return "foos"
}

func TestSave_hooks(t *testing.T) {
	test.Run(t, "runs hooks", func(t *testing.T, tx *sqlx.Tx) {
		f := &FooSaveHookTest{
			Foo: test.Foo{
				ID: 1,
			},
		}
		err := model.Save(tx, f)
		assert.NoError(t, err)

		assert.True(t, f.saved)
	})

	test.Run(t, "runs hooks on anonymise structs", func(t *testing.T, tx *sqlx.Tx) {
		f := &FooSaveHookTestWrapper{
			FooSaveHookTest{
				Foo: test.Foo{
					ID: 1,
				},
			},
		}
		err := model.Save(tx, f)
		assert.NoError(t, err)

		assert.True(t, f.saved)
	})
}

type SaveFooReadonly struct {
	test.Foo
	Readonly string `db:"computed,readonly"`
}

func (f *SaveFooReadonly) Table() string {
	return "foos"
}

func TestSave_readonly(t *testing.T) {
	test.Run(t, "runs hooks", func(t *testing.T, tx *sqlx.Tx) {
		f := &SaveFooReadonly{
			Foo: test.Foo{
				ID: 1,
			},
			Readonly: "yes",
		}
		err := model.Save(tx, f)
		assert.NoError(t, err)

		newFoo, err := builder.From[*SaveFooReadonly]().First(tx)
		assert.NoError(t, err)
		assert.Equal(t, "", newFoo.Readonly)
	})

}

func TestInsertManyContext(t *testing.T) {
	test.Run(t, "insert", func(t *testing.T, tx *sqlx.Tx) {
		err := model.InsertManyContext(context.TODO(), tx, []*test.Foo{{Name: "1"}, {Name: "2"}, {Name: "3"}})
		assert.NoError(t, err)

		foos, err := builder.From[*SaveFooReadonly]().Get(tx)
		assert.NoError(t, err)
		assert.Len(t, foos, 3)
		for i, foo := range foos {
			assert.True(t, foo.InDatabase())
			assert.Equal(t, fmt.Sprint(i+1), foo.Name)
		}
	})

}
