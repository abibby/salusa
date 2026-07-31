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
	"github.com/abibby/salusa/database/model/mixins"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

type Case struct {
	Name             string
	Builder          dialects.QueryBuilder
	ExpectedSQL      string
	ExpectedBindings []any
}

func QueryTest(t *testing.T, testCases []Case) {
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			result, err := dialects.New().EncodeSelectQuery(tc.Builder.Query())
			if assert.NoError(t, err) {
				assert.Equal(t, tc.ExpectedSQL, result.Query)
				assert.Equal(t, tc.ExpectedBindings, result.Bindings)
			}
		})
	}
}

type EncoderTestCase[T any] struct {
	Name             string
	Builder          T
	ExpectedSQL      string
	ExpectedBindings []any
	ExpectedError    error
}

func EncoderTest[T any](t *testing.T, encoder func(v T) (dialects.SQLResult, error), testCases []EncoderTestCase[T]) {
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			result, err := encoder(tc.Builder)
			if assert.NoError(t, err) {
				assert.Equal(t, tc.ExpectedSQL, result.Query)
				assert.Equal(t, tc.ExpectedBindings, result.Bindings)
			}
		})
	}
}

var runner = dbtest.NewRunner(func() (*sqlx.DB, error) {
	cfg := sqlite.NewConfig(":memory:")
	cfg.SetDialect()
	db, err := sqlx.Open(cfg.DriverName(), cfg.DataSourceName())
	if err != nil {
		return nil, err
	}
	ctx := context.Background()
	err = migrate.RunModelCreate(ctx, db, &Foo{}, &Bar{}, &FooSoftDelete{})
	if err != nil {
		return nil, err
	}
	return db, nil
})

var Run = runner.Run
var RunNoTx = runner.RunNoTx
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

type FooSoftDelete struct {
	model.BaseModel
	mixins.SoftDelete
	ID   int    `json:"id"   db:"id,primary,autoincrement"`
	Name string `json:"name" db:"name"`
}

func (h *FooSoftDelete) Table() string {
	return "foo_soft_deletes"
}
