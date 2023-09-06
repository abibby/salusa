package builder

import (
	"github.com/abibby/salusa/database/dialects"
	"github.com/abibby/salusa/internal/helpers"
)

type groupBys []helpers.ToSQLer

func (g groupBys) Clone() groupBys {
	return cloneSlice(g)
}
func (g groupBys) ToSQL(d dialects.Dialect) (string, []any, error) {
	if len(g) == 0 {
		return "", nil, nil
	}
	r := helpers.Result()
	r.AddString("GROUP BY")
	r.Add(helpers.Join(g, ", "))
	return r.ToSQL(d)
}

// GroupBy sets the "group by" clause to the query.
func (b groupBys) GroupBy(columns ...string) groupBys {
	return helpers.IdentifierList(columns)
}

// GroupBy adds a "group by" clause to the query.
func (b groupBys) AddGroupBy(columns ...string) groupBys {
	return append(b, helpers.IdentifierList(columns)...)
}
