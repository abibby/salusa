package helpers

import (
	"github.com/abibby/salusa/database/dialects"
)

type ToSQLer interface {
	ToSQL(d dialects.Dialect) (string, []any, error)
}

type ToSQLFunc func(d dialects.Dialect) (string, []any, error)

func (f ToSQLFunc) ToSQL(d dialects.Dialect) (string, []any, error) {
	return f(d)
}

func Identifier(i string) ToSQLer {
	return ToSQLFunc(func(d dialects.Dialect) (string, []any, error) {
		return d.Identifier(i), nil, nil
	})
}

func IdentifierList(strs []string) []ToSQLer {
	identifiers := make([]ToSQLer, len(strs))
	for i, s := range strs {
		identifiers[i] = Identifier(s)
	}
	return identifiers
}

func Join[T ToSQLer](sqlers []T, sep string) ToSQLer {
	return ToSQLFunc(func(d dialects.Dialect) (string, []any, error) {
		if sqlers == nil {
			return "", []any{}, nil
		}
		sql := ""
		bindings := []any{}
		for i, sqler := range sqlers {
			q, b, err := sqler.ToSQL(d)
			if err != nil {
				return "", nil, err
			}
			sql += q
			if i < len(sqlers)-1 {
				sql += sep
			}
			bindings = append(bindings, b...)
		}
		return sql, bindings, nil
	})
}

func Raw(sql string, bindings ...any) ToSQLer {
	return ToSQLFunc(func(d dialects.Dialect) (string, []any, error) {
		return sql, bindings, nil
	})
}

func Group(sqler ToSQLer) ToSQLer {
	return ToSQLFunc(func(d dialects.Dialect) (string, []any, error) {
		q, bindings, err := sqler.ToSQL(d)
		return "(" + q + ")", bindings, err
	})
}
func Concat(sqlers ...ToSQLer) ToSQLer {
	return Join(sqlers, "")
}

func Literal(v any) ToSQLer {
	return ToSQLFunc(func(d dialects.Dialect) (string, []any, error) {
		return d.Binding(), []any{v}, nil
	})
}

func LiteralList(values []any) []ToSQLer {
	literals := make([]ToSQLer, len(values))
	for i, s := range values {
		literals[i] = Literal(s)
	}
	return literals
}
