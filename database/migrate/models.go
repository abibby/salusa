package migrate

import (
	"context"
	"fmt"

	"github.com/abibby/salusa/database/internal/helpers"
	"github.com/abibby/salusa/database/models"
)

func RunModelCreate(ctx context.Context, db helpers.QueryExecer, models ...models.Model) error {
	for _, model := range models {
		m, err := CreateFromModel(model)
		if err != nil {
			return fmt.Errorf("migration for %s: %w", helpers.GetTable(model), err)
		}
		err = m.Run(ctx, db)
		if err != nil {
			return fmt.Errorf("migration for %s: %w", helpers.GetTable(model), err)
		}
	}
	return nil
}
