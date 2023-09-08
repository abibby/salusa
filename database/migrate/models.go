package migrate

import (
	"context"
	"fmt"

	"github.com/abibby/salusa/database/model"
	"github.com/abibby/salusa/internal/helpers"
)

func RunModelCreate(ctx context.Context, db helpers.QueryExecer, models ...model.Model) error {
	for _, m := range models {
		m, err := CreateFromModel(m)
		if err != nil {
			return fmt.Errorf("migration for %s: %w", helpers.GetTable(m), err)
		}
		err = m.Run(ctx, db)
		if err != nil {
			return fmt.Errorf("migration for %s: %w", helpers.GetTable(m), err)
		}
	}
	return nil
}
