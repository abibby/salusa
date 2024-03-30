package builder_test

import (
	"fmt"

	"github.com/abibby/salusa/database/builder"
	"github.com/abibby/salusa/database/dialects"
	"github.com/abibby/salusa/internal/test"
)

func ExampleBuilder() {
	query, bindings, err := builder.
		From[*test.Foo]().
		Where("column", "=", "value").
		ToSQL(dialects.New())
	if err != nil {
		panic(err)
	}

	fmt.Println(bindings)
	fmt.Println(query)
	// Output:
	// [value]
	// SELECT "foos".* FROM "foos" WHERE "column" = ?
}

func ExampleBuilder_WhereHas() {
	query, bindings, err := builder.
		From[*test.Foo]().
		WhereHas("Bar", func(q *builder.Builder) *builder.Builder {
			return q.Where("id", "=", 7)
		}).
		ToSQL(dialects.New())
	if err != nil {
		panic(err)
	}

	fmt.Println(bindings)
	fmt.Println(query)
	// Output:
	// [7]
	// SELECT "foos".* FROM "foos" WHERE EXISTS (SELECT "bars".* FROM "bars" WHERE "foo_id" = "foos"."id" AND "id" = ?)
}
