package schema

import (
	"context"

	"github.com/abibby/salusa/database"
	"github.com/abibby/salusa/internal/helpers"
)

func Drop(table string) Runner {
	return Run(func(ctx context.Context, tx database.DB) error {
		return runQuery(ctx, tx, helpers.Concat(helpers.Raw("DROP TABLE "), helpers.Identifier(table)))
	})
}
func DropIfExists(table string) Runner {
	return Run(func(ctx context.Context, tx database.DB) error {
		return runQuery(ctx, tx, helpers.Concat(helpers.Raw("DROP TABLE IF EXISTS "), helpers.Identifier(table)))
	})
}
