package schema

import (
	"context"
	"fmt"

	"gosalusa.com/database"
)

type ViewBuilder struct {
	Name  string
	Query string
}

func View(name string, query string) *ViewBuilder {
	return &ViewBuilder{
		Name:  name,
		Query: query,
	}
}

func (b *ViewBuilder) GoString() string {
	return fmt.Sprintf(
		"schema.View(%#v, %#v)",
		b.Name,
		b.Query,
	)
}
func (b *ViewBuilder) Run(ctx context.Context, tx database.DB) error {
	_, err := tx.ExecContext(ctx, fmt.Sprintf("CREATE VIEW %s as %s"))
	return err
}
