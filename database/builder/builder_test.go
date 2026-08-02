package builder_test

import (
	"fmt"

	"github.com/abibby/salusa/database/builder"
	"github.com/abibby/salusa/database/dialects/sqlite"
	"github.com/abibby/salusa/internal/test"
)

func ExampleBuilder() {
	q := builder.
		From[*test.Foo]().
		Where("column", "=", "value")

	r, err := sqlite.New().EncodeSelectQuery(q.Query())
	if err != nil {
		panic(err)
	}

	fmt.Println(r.SQL)
	fmt.Println(r.Bindings)
	// Output:
	// SELECT "foos".* FROM "foos" WHERE "column" = ?
	// [value]
}

func ExampleBuilder_WhereHas() {
	q := builder.
		From[*test.Foo]().
		WhereHas("Bar", func(q *builder.Builder) *builder.Builder {
			return q.Where("id", "=", 7)
		})
	r, err := sqlite.New().EncodeSelectQuery(q.Query())
	if err != nil {
		panic(err)
	}

	fmt.Println(r.SQL)
	fmt.Println(r.Bindings)
	// Output:
	// SELECT "foos".* FROM "foos" WHERE EXISTS (SELECT "bars".* FROM "bars" WHERE "foo_id" = "foos"."id" AND "id" = ?)
	// [7]
}
