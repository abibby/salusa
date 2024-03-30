package builder

import (
	"github.com/abibby/salusa/database"
	"github.com/abibby/salusa/database/dialects"
	"github.com/abibby/salusa/internal/helpers"
)

type Deleter struct {
	builder *Builder
}

func (d *Deleter) ToSQL(dialect dialects.Dialect) (string, []any, error) {
	return helpers.Concat(
		helpers.Raw("DELETE "),
		d.builder.Select(),
	).ToSQL(dialect)
}

func (b *ModelBuilder[T]) Delete(tx database.DB) error {
	return b.subBuilder.Delete(tx)
}
func (b *Builder) Delete(tx database.DB) error {
	q, bindings, err := b.Deleter().ToSQL(dialects.New())
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(b.ctx, q, bindings...)
	if err != nil {
		return err
	}

	return nil
}

func (b *ModelBuilder[T]) Deleter() *Deleter {
	return b.subBuilder.Deleter()
}
func (b *Builder) Deleter() *Deleter {
	return &Deleter{
		builder: b,
	}
}
