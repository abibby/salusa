package builder

import (
	"github.com/abibby/salusa/database/dialects"
	"github.com/abibby/salusa/database/internal/helpers"
)

type limit struct {
	limit  int
	offset int
}

func (l *limit) Clone() *limit {
	return &limit{
		limit:  l.limit,
		offset: l.offset,
	}
}
func (l *limit) ToSQL(d dialects.Dialect) (string, []any, error) {
	if l.limit == 0 && l.offset == 0 {
		return "", nil, nil
	}
	r := helpers.Result()
	r.AddString("LIMIT")
	r.Add(helpers.Literal(l.limit))
	if l.offset != 0 {
		r.AddString("OFFSET")
		r.Add(helpers.Literal(l.offset))
	}
	return r.ToSQL(d)
}

// Limit set the maximum number of rows to return.
func (l *limit) Limit(limit int) *limit {
	l.limit = limit
	return l
}

// Offset sets the number of rows to skip before returning the result.
func (l *limit) Offset(offset int) *limit {
	l.offset = offset
	return l
}
