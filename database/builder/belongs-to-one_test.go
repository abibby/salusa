package builder_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/abibby/salusa/database/builder"
	"github.com/abibby/salusa/database/dialects/sqlite"
	"github.com/abibby/salusa/database/insert"
	"github.com/abibby/salusa/database/migrate"
	"github.com/abibby/salusa/database/models"
	"github.com/abibby/salusa/internal/test"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func ExampleBelongsTo() {
	sqlite.UseSQLite()
	type Bar struct {
		models.BaseModel
		ID   int    `db:"id,autoincrement,primary"`
		Name string `db:"name"`
	}

	type Foo struct {
		models.BaseModel
		ID    int `db:"id,autoincrement,primary"`
		BarID int `db:"bar_id"`
		Bar   *builder.BelongsTo[*Bar]
	}

	db, _ := sqlx.Open("sqlite3", ":memory:")

	createFoo, _ := migrate.CreateFromModel(&Foo{})
	createFoo.Run(context.Background(), db)
	createBar, _ := migrate.CreateFromModel(&Bar{})
	createBar.Run(context.Background(), db)

	foo := &Foo{BarID: 1}
	insert.Save(db, foo)
	bar := &Bar{ID: 1, Name: "bar name"}
	insert.Save(db, bar)

	builder.Load(db, foo, "Bar")
	relatedBar, _ := foo.Bar.Value()

	fmt.Println(relatedBar.Name)

	// Output: bar name
}

func TestBelongsToLoad(t *testing.T) {
	test.Run(t, "", func(t *testing.T, tx *sqlx.Tx) {
		foos := []*test.Foo{
			{ID: 1},
			{ID: 2},
			{ID: 3},
		}
		for _, f := range foos {
			assert.NoError(t, insert.Save(tx, f))
		}
		bars := []*test.Bar{
			{ID: 4, FooID: 1},
			{ID: 5, FooID: 2},
			{ID: 6, FooID: 3},
		}
		for _, b := range bars {
			assert.NoError(t, insert.Save(tx, b))
		}

		err := builder.Load(tx, bars, "Foo")
		if !assert.NoError(t, err) {
			return
		}

		for _, bar := range bars {
			assert.True(t, bar.Foo.Loaded())
			foo, ok := bar.Foo.Value()
			assert.True(t, ok)
			assert.Equal(t, bar.FooID, foo.ID)
		}
	})

	test.Run(t, "uuids", func(t *testing.T, tx *sqlx.Tx) {
		type Foo struct {
			models.BaseModel
			ID   int       `json:"id" db:"id,primary,autoincrement"`
			Name uuid.UUID `json:"name" db:"name"`
		}
		type Bar struct {
			models.BaseModel
			FooName *uuid.UUID               `json:"foo_id" db:"foo_id"`
			Foo     *builder.BelongsTo[*Foo] `json:"foo"    db:"-" foreign:"foo_id" owner:"name"`
		}

		bars := []*Bar{
			{FooName: newUUID()},
			{FooName: newUUID()},
			{FooName: nil},
			{FooName: nil},
		}
		for _, b := range bars {
			if b.FooName != nil {
				MustSave(tx, &Foo{Name: *b.FooName})
			}
		}
		builder.InitializeRelationships(bars)
		builder.Load(tx, bars, "Foo")

		for i, b := range bars {
			f, ok := b.Foo.Value()
			assert.True(t, ok)
			if i < 2 {
				assert.NotNil(t, f)
			} else {
				assert.Nil(t, f)
			}
		}
	})
}

func newUUID() *uuid.UUID {
	id := uuid.New()
	return &id
}
