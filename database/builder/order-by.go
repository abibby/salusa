package builder

import (
	"github.com/abibby/salusa/database/dialects"
	"github.com/abibby/salusa/internal/helpers"
)

type orderBys []helpers.SQLStringer

func (o orderBys) Clone() orderBys {
	return cloneSlice(o)
}
func (o orderBys) SQLString(d dialects.Dialect) (string, []any, error) {
	if len(o) == 0 {
		return "", nil, nil
	}
	r := helpers.Result()
	r.AddString("ORDER BY")
	r.Add(helpers.Join(o, ", "))
	return r.SQLString(d)
}

// OrderBy adds an order by clause to the query.
func (o orderBys) OrderBy(column string) orderBys {
	return append(o, helpers.Identifier(column))
}

// OrderByDesc adds a descending order by clause to the query.
func (o orderBys) OrderByDesc(column string) orderBys {
	return append(o, helpers.Join([]helpers.SQLStringer{helpers.Identifier(column), helpers.Raw("DESC")}, " "))
}

// Unordered removes all order by clauses from the query.
func (o orderBys) Unordered() orderBys {
	return orderBys{}
}
