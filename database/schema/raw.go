package schema

import (
	"context"

	"abibby.com/salusa/database"
)

type Raw string

// Run implements [Runner].
func (v Raw) Run(ctx context.Context, tx database.DB) error {
	_, err := tx.ExecContext(ctx, string(v))
	return err
}
