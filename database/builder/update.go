package builder

import (
	"errors"

	"github.com/abibby/salusa/database"
	"github.com/abibby/salusa/database/dialects"
	"github.com/abibby/salusa/internal/helpers"
)

var (
	ErrNoUpdates = errors.New("no updates found")
)

type Updates map[string]any
type Updater struct {
	builder *Builder
	updates Updates
}

func (d *Updater) SQLString(dialect dialects.Dialect) (string, []any, error) {
	if len(d.updates) == 0 {
		return "", nil, ErrNoUpdates
	}

	sets := []helpers.SQLStringer{}

	for k, v := range d.updates {
		sets = append(sets, helpers.Concat(
			helpers.Identifier(k),
			helpers.Raw("="),
			helpers.Literal(v),
		))
	}

	parts := []helpers.SQLStringer{
		helpers.Raw("UPDATE"),
		helpers.Identifier(string(d.builder.from)),
		helpers.Raw("SET"),
		helpers.Join(sets, ", "),
	}

	if len(d.builder.wheres.list) > 0 {
		parts = append(parts, d.builder.wheres)
	}

	return helpers.Join(parts, " ").SQLString(dialect)
}

func (b *ModelBuilder[T]) Update(tx database.DB, updates Updates) error {
	return b.builder.Update(tx, updates)
}
func (b *Builder) Update(tx database.DB, updates Updates) error {
	if len(updates) == 0 {
		return nil
	}

	q, bindings, err := b.Updater(updates).SQLString(dialects.New())
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(b.ctx, q, bindings...)
	if err != nil {
		return err
	}

	return nil
}

func (b *ModelBuilder[T]) Updater(updates Updates) *Updater {
	return b.builder.Updater(updates)
}
func (b *Builder) Updater(updates Updates) *Updater {
	return &Updater{
		builder: b,
		updates: updates,
	}
}
