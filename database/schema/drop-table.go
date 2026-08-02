package schema

import (
	"context"

	"github.com/abibby/salusa/database"
	"github.com/abibby/salusa/database/dialects"
)

func Drop(table string) Runner {
	return Run(func(ctx context.Context, tx database.DB) error {
		return runDropTable(ctx, tx, &dialects.DropTableQuery{Table: table})
	})
}
func DropIfExists(table string) Runner {
	return Run(func(ctx context.Context, tx database.DB) error {
		return runDropTable(ctx, tx, &dialects.DropTableQuery{Table: table, IfExists: true})
	})
}
func runDropTable(ctx context.Context, tx database.DB, q *dialects.DropTableQuery) error {
	result, err := dialects.New().EncodeDropTableQuery(q)
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, result.SQL, result.Bindings...)
	return err
}
