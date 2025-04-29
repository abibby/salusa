package builder

import (
	"fmt"
	"reflect"

	"github.com/abibby/salusa/database"
	"github.com/abibby/salusa/database/dialects"
	"github.com/abibby/salusa/database/hooks"
	"github.com/abibby/salusa/internal/helpers"
	"github.com/abibby/salusa/internal/relationship"
	"github.com/jmoiron/sqlx"
)

type QueryError struct {
	err   error
	query string
}

func (e *QueryError) Error() string {
	return fmt.Sprintf("%s: %v", e.query, e.err)
}

func (e *QueryError) Unwrap() error {
	return e.err
}

// Get executes the query as a select statement and returns the result.
func (b *ModelBuilder[T]) Get(tx database.DB) ([]T, error) {
	v := []T{}
	err := b.builder.Load(tx, &v)
	if err != nil {
		return nil, err
	}

	for _, with := range b.withs {
		err = LoadMissingContext(b.Context(), tx, v, with)
		if err != nil {
			return nil, err
		}
	}

	return v, nil
}

// Get executes the query as a select statement and returns the first record.
func (b *ModelBuilder[T]) First(tx database.DB) (T, error) {
	v, err := b.Clone().
		Limit(1).
		Get(tx)
	if err != nil {
		var zero T
		return zero, err
	}
	if len(v) < 1 {
		var zero T
		return zero, nil
	}
	return v[0], nil
}

// Find returns the record with a matching primary key. It will fail on tables with multiple primary keys.
func (b *ModelBuilder[T]) Find(tx database.DB, primaryKeyValue any) (T, error) {
	var m T
	pKeys := helpers.PrimaryKey(m)
	if len(pKeys) != 1 {
		return m, fmt.Errorf("Find only supports tables with 1 primary key")
	}
	return b.
		Where(pKeys[0], "=", primaryKeyValue).
		First(tx)
}

// Load executes the query as a select statement and sets v to the result.
func (b *ModelBuilder[T]) Load(tx database.DB, v any) error {
	return b.builder.Load(tx, v)
}

// Load executes the query as a select statement and sets v to the result.
func (b *Builder) Load(tx database.DB, v any) (err error) {
	q, bindings, err := b.SQLString(dialects.New())
	if err != nil {
		return err
	}
	defer func() {
		if err == nil {
			return
		}
		err = &QueryError{
			err:   err,
			query: q,
		}
	}()
	if reflect.TypeOf(v).Elem().Kind() == reflect.Slice {
		err = sqlx.SelectContext(b.Context(), tx, v, q, bindings...)
	} else {
		err = sqlx.GetContext(b.Context(), tx, v, q, bindings...)
	}
	if err != nil {
		return err
	}

	err = relationship.InitializeRelationships(v)
	if err != nil {
		return err
	}

	err = hooks.AfterLoad(b.Context(), tx, v)
	if err != nil {
		return err
	}
	return nil
}

// Load executes the query as a select statement and sets v to the result.
//
// Deprecated: Use ModelBuilder.Load
func (b *ModelBuilder[T]) LoadOne(tx database.DB, v any) error {
	return b.builder.LoadOne(tx, v)
}

// Load executes the query as a select statement and sets v to the first record.
//
// Deprecated: Use Builder.Load
func (b *Builder) LoadOne(tx database.DB, v any) error {
	return b.Load(tx, v)
}

func (b *ModelBuilder[T]) Each(tx database.DB, cb func(v T) error) error {
	return b.Chunk(tx, func(models []T) error {
		for _, model := range models {
			err := cb(model)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (b *ModelBuilder[T]) Chunk(tx database.DB, cb func(v []T) error) error {
	return b.ChunkN(tx, 1000, cb)
}
func (b *ModelBuilder[T]) ChunkN(tx database.DB, limit int, cb func(v []T) error) error {
	offset := 0
	for {
		models, err := b.Limit(limit).Offset(offset).Get(tx)
		if err != nil {
			return err
		}
		offset += limit
		if len(models) == 0 {
			return nil
		}
		err = cb(models)
		if err != nil {
			return err
		}
	}
}

// Count executes select and returns the number of records.
func (b *ModelBuilder[T]) Count(tx database.DB) (int, error) {
	return b.builder.Count(tx)
}

// Count executes select and returns the number of records.
func (b *Builder) Count(tx database.DB) (int, error) {
	return b.numericFunc(tx, "count", "*")
}

func (b *Builder) numericFunc(tx database.DB, function, column string) (int, error) {
	var count int
	err := b.
		Unordered().
		SelectFunction(function, column).LoadOne(tx, &count)
	if err != nil {
		return 0, err
	}
	return count, nil
}
