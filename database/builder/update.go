package builder

import (
	"errors"

	"github.com/abibby/salusa/database"
	"github.com/abibby/salusa/database/dialects"
)

var (
	ErrNoUpdates = errors.New("no updates found")
)

type Updates map[string]any

func (b *ModelBuilder[T]) Update(tx database.DB, updates Updates) error {
	return b.builder.Update(tx, updates)
}
func (b *Builder) Update(tx database.DB, updates Updates) error {
	if len(updates) == 0 {
		return nil
	}
	r, err := dialects.New().EncodeUpdateQuery(&dialects.UpdateQuery{
		Table:  b.GetTable(),
		Values: updates,
		Wheres: b.wheres.conditions,
	})
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(b.ctx, r.SQL, r.Bindings...)
	if err != nil {
		return err
	}

	return nil
}
