package builder

import (
	"github.com/abibby/salusa/database"
	"github.com/abibby/salusa/database/dialects"
	"github.com/abibby/salusa/internal/helpers"
)

type Deleter struct {
	builder *Builder
}

func (d *Deleter) SQLString(dialect dialects.Dialect) (string, []any, error) {
	return helpers.Concat(
		helpers.Raw("DELETE "),
		d.builder.Select(),
	).SQLString(dialect)
}

func (b *ModelBuilder[T]) Delete(tx database.DB) error {
	return b.builder.Delete(tx)
}
func (b *Builder) Delete(tx database.DB) error {
	q, bindings, err := b.Deleter().SQLString(dialects.New())
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
	return b.builder.Deleter()
}
func (b *Builder) Deleter() *Deleter {
	return &Deleter{
		builder: b,
	}
}
