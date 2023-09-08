package test

import (
	"context"
	"testing"

	"github.com/abibby/salusa/database/builder"
	"github.com/abibby/salusa/database/dbtest"
	"github.com/abibby/salusa/database/dialects"
	"github.com/abibby/salusa/database/dialects/sqlite"
	"github.com/abibby/salusa/database/migrate"
	"github.com/abibby/salusa/database/model"
	"github.com/abibby/salusa/internal/helpers"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

type Case struct {
	Name             string
	Builder          helpers.ToSQLer
	ExpectedSQL      string
	ExpectedBindings []any
}

func QueryTest(t *testing.T, testCases []Case) {

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			q, bindings, err := tc.Builder.ToSQL(dialects.New())
			assert.NoError(t, err)

			assert.Equal(t, tc.ExpectedSQL, q)
			assert.Equal(t, tc.ExpectedBindings, bindings)
		})
	}
}

var runner = dbtest.NewRunner(func() (*sqlx.DB, error) {
	sqlite.UseSQLite()
	db, err := sqlx.Open("sqlite3", ":memory:")
	if err != nil {
		return nil, err
	}
	ctx := context.Background()
	err = migrate.RunModelCreate(ctx, db, &Foo{}, &Bar{})
	if err != nil {
		return nil, err
	}
	return db, nil
})

var Run = runner.Run
var RunBenchmark = runner.RunBenchmark

type Foo struct {
	model.BaseModel
	ID   int                    `json:"id"   db:"id,primary,autoincrement"`
	Name string                 `json:"name" db:"name"`
	Bar  *builder.HasOne[*Bar]  `json:"bar"`
	Bars *builder.HasMany[*Bar] `json:"bars"`
}

func (h *Foo) Table() string {
	return "foos"
}

type Bar struct {
	model.BaseModel
	ID    int                      `json:"id"     db:"id,primary,autoincrement"`
	FooID int                      `json:"foo_id" db:"foo_id"`
	Foo   *builder.BelongsTo[*Foo] `json:"foo"`
}

func (h *Bar) Table() string {
	return "bars"
}
