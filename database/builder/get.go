package builder

import (
	"fmt"

	"github.com/abibby/salusa/database"
	"github.com/abibby/salusa/database/dialects"
	"github.com/abibby/salusa/database/hooks"
	"github.com/abibby/salusa/internal/helpers"
	"github.com/abibby/salusa/internal/relationship"
	"github.com/jmoiron/sqlx"
)

// Get executes the query as a select statement and returns the result.
func (b *Builder[T]) Get(tx database.DB) ([]T, error) {
	v := []T{}
	err := b.Load(tx, &v)
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
func (b *Builder[T]) First(tx database.DB) (T, error) {
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
	return v[0], err
}

// Find returns the record with a matching primary key. It will fail on tables with multiple primary keys.
func (b *Builder[T]) Find(tx database.DB, primaryKeyValue any) (T, error) {
	var m T
	pKeys := helpers.PrimaryKey(m)
	if len(pKeys) != 1 {
		return m, fmt.Errorf("Find only supports tables with 1 primary key")
	}
	return b.Clone().
		Where(pKeys[0], "=", primaryKeyValue).
		First(tx)
}

// Load executes the query as a select statement and sets v to the result.
func (b *Builder[T]) Load(tx database.DB, v any) error {
	q, bindings, err := b.ToSQL(dialects.New())
	if err != nil {
		return err
	}

	err = sqlx.SelectContext(b.Context(), tx, v, q, bindings...)
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

// Load executes the query as a select statement and sets v to the first record.
func (b *Builder[T]) LoadOne(tx database.DB, v any) error {
	q, bindings, err := b.Clone().
		Limit(1).
		ToSQL(dialects.New())

	if err != nil {
		return err
	}

	err = sqlx.GetContext(b.Context(), tx, v, q, bindings...)
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

func (b *Builder[T]) Each(tx database.DB, cb func(v T) error) error {
	limit := 1000
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
		for _, model := range models {
			err = cb(model)
			if err != nil {
				return err
			}
		}
	}
}
