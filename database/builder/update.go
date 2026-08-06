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
	d, err := dialects.New(tx.DriverName())
	if err != nil {
		return err
	}
	r, err := d.EncodeUpdateQuery(b.UpdateQuery(updates))
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(b.ctx, r.SQL, r.Bindings...)
	if err != nil {
		return err
	}

	return nil
}

func (b *ModelBuilder[T]) UpdateQuery(updates Updates) *dialects.UpdateQuery {
	return b.builder.UpdateQuery(updates)
}

func (b *Builder) UpdateQuery(updates Updates) *dialects.UpdateQuery {
	return &dialects.UpdateQuery{
		Table:  b.GetTable(),
		Values: updates,
		Wheres: b.wheres.conditions,
	}
}
