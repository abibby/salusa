package dialects

import (
	"errors"
	"fmt"
)

var ErrNotRegistered = errors.New("no dialect registered")

type Dialect interface {
	EncodeSelectQuery(q *SelectQuery) (RawQuery, error)
	EncodeInsertQuery(q *InsertQuery) (RawQuery, error)
	EncodeUpdateQuery(q *UpdateQuery) (RawQuery, error)
	EncodeDeleteQuery(q *DeleteQuery) (RawQuery, error)
	EncodeCreateTableQuery(q *CreateTableQuery) (RawQuery, error)
	EncodeDropTableQuery(q *DropTableQuery) (RawQuery, error)
	EncodeAlterTableQuery(q *AlterTableQuery) (RawQuery, error)
}

var dialects = map[string]func() Dialect{}

func Register(driver string, factory func() Dialect) {
	dialects[driver] = factory
}

func New(driverName string) (Dialect, error) {
	f, ok := dialects[driverName]
	if !ok {
		return nil, fmt.Errorf("%w for %s", ErrNotRegistered, driverName)
	}
	return f(), nil
}
