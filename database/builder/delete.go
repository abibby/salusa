package builder

import (
	"github.com/abibby/salusa/database/dialects"
	"github.com/abibby/salusa/internal/helpers"
)

type Deleter struct {
	builder *SubBuilder
}

func (d *Deleter) ToSQL(dialect dialects.Dialect) (string, []any, error) {
	return helpers.Concat(
		helpers.Raw("DELETE "),
		d.builder.Select(),
	).ToSQL(dialect)
}

func (b *Builder[T]) Delete(tx helpers.QueryExecer) error {
	return b.subBuilder.Delete(tx)
}
func (b *SubBuilder) Delete(tx helpers.QueryExecer) error {
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

func (b *Builder[T]) Deleter() *Deleter {
	return b.subBuilder.Deleter()
}
func (b *SubBuilder) Deleter() *Deleter {
	return &Deleter{
		builder: b,
	}
}
