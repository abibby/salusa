package helpers

import (
	"github.com/abibby/salusa/database/dialects"
)

type SQLStringer interface {
	SQLString(d dialects.Dialect) (string, []any, error)
}

type SQLStringFunc func(d dialects.Dialect) (string, []any, error)

func (f SQLStringFunc) SQLString(d dialects.Dialect) (string, []any, error) {
	return f(d)
}

func Identifier(i string) SQLStringer {
	return SQLStringFunc(func(d dialects.Dialect) (string, []any, error) {
		return d.Identifier(i), nil, nil
	})
}

func IdentifierList(strs []string) []SQLStringer {
	identifiers := make([]SQLStringer, len(strs))
	for i, s := range strs {
		identifiers[i] = Identifier(s)
	}
	return identifiers
}

func Join[T SQLStringer](sqlers []T, sep string) SQLStringer {
	return SQLStringFunc(func(d dialects.Dialect) (string, []any, error) {
		if sqlers == nil {
			return "", []any{}, nil
		}
		sql := ""
		bindings := []any{}
		for i, sqler := range sqlers {
			q, b, err := sqler.SQLString(d)
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

func Raw(sql string, bindings ...any) SQLStringer {
	return SQLStringFunc(func(d dialects.Dialect) (string, []any, error) {
		return sql, bindings, nil
	})
}

func Group(sqler SQLStringer) SQLStringer {
	return SQLStringFunc(func(d dialects.Dialect) (string, []any, error) {
		q, bindings, err := sqler.SQLString(d)
		return "(" + q + ")", bindings, err
	})
}
func Concat(sqlers ...SQLStringer) SQLStringer {
	return Join(sqlers, "")
}

func Literal(v any) SQLStringer {
	return SQLStringFunc(func(d dialects.Dialect) (string, []any, error) {
		return d.Binding(), []any{v}, nil
	})
}

func LiteralList(values []any) []SQLStringer {
	literals := make([]SQLStringer, len(values))
	for i, s := range values {
		literals[i] = Literal(s)
	}
	return literals
}
