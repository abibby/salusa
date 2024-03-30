package helpers

import (
	"github.com/abibby/salusa/database/dialects"
)

type SQLResult struct {
	sqlers []SQLStringer
}

func Result() *SQLResult {
	return &SQLResult{}
}

func (r *SQLResult) AddString(sql string) *SQLResult {
	return r.Add(Raw(sql))
}
func (r *SQLResult) Add(sqler SQLStringer) *SQLResult {
	r.sqlers = append(r.sqlers, sqler)
	return r
}

func (r *SQLResult) SQLString(d dialects.Dialect) (string, []any, error) {
	resultSql := ""
	resultBindings := []any{}
	for _, sqler := range r.sqlers {
		sql, bindings, err := sqler.SQLString(d)
		if err != nil {
			return "", nil, err
		}
		if resultSql != "" && sql != "" {
			resultSql += " "
		}
		resultSql += sql

		if bindings != nil {
			resultBindings = append(resultBindings, bindings...)
		}

	}
	return resultSql, resultBindings, nil
}
