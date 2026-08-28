package builder

import (
	"errors"

	"abibby.com/salusa/database"
	"abibby.com/salusa/database/dialects"
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

func (b *ModelBuilder[T]) UpdateReturning(tx database.DB, updates Updates) ([]T, error) {
	if len(updates) == 0 {
		return nil, nil
	}
	d, err := dialects.New(tx.DriverName())
	if err != nil {
		return nil, err
	}

	// b.Query().Select.Columns
	q := b.UpdateQuery(updates)
	q.Returning = []dialects.Column{
		{Column: "*"},
	}

	r, err := d.EncodeUpdateQuery(q)
	if err != nil {
		return nil, err
	}

	result := []T{}

	err = load(b.Context(), tx, r, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
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
