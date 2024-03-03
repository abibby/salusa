package migrate

import (
	"context"
	"fmt"

	"github.com/abibby/salusa/database"
	"github.com/abibby/salusa/database/model"
)

func RunModelCreate(ctx context.Context, db database.DB, models ...model.Model) error {
	for _, m := range models {
		m, err := CreateFromModel(m)
		if err != nil {
			return fmt.Errorf("migration for %s: %w", database.GetTable(m), err)
		}
		err = m.Run(ctx, db)
		if err != nil {
			return fmt.Errorf("migration for %s: %w", database.GetTable(m), err)
		}
	}
	return nil
}
