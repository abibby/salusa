package builder

import (
	"gosalusa.com/database"
	"gosalusa.com/database/dialects"
)

func (b *ModelBuilder[T]) Delete(tx database.DB) error {
	return b.builder.Delete(tx)
}

func delete(b *Builder, tx database.DB) error {
	d, err := dialects.New(tx.DriverName())
	if err != nil {
		return err
	}
	result, err := d.EncodeDeleteQuery(b.DeleteQuery())
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(b.ctx, result.SQL, result.Bindings...)
	if err != nil {
		return err
	}
	return nil
}
func (b *Builder) Delete(tx database.DB) error {
	current := delete
	for _, s := range b.ActiveScopes() {
		if s.Delete != nil {
			current = s.Delete(current)
		}
	}
	return current(b, tx)
}

func (b *ModelBuilder[T]) DeleteQuery() *dialects.DeleteQuery {
	return b.builder.DeleteQuery()
}
func (b *Builder) DeleteQuery() *dialects.DeleteQuery {
	return &dialects.DeleteQuery{
		Table:  b.GetTable(),
		Wheres: b.wheres.conditions,
	}
}
