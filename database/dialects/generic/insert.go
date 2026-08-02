package generic

import (
	"errors"
	"strings"

	"github.com/abibby/salusa/database/dialects"
	"github.com/abibby/salusa/slices"
)

var ErrInsertNoRows = errors.New("no rows to insert")
var ErrInsertMismatchedValueKeys = errors.New("mismatched value keys")

func (g *Generic) EncodeInsertQuery(q *dialects.InsertQuery) (dialects.RawQuery, error) {
	if len(q.Values) == 0 {
		return dialects.RawQuery{}, ErrInsertNoRows
	}
	numColumns := len(q.Values[0])

	columns := make([]string, 0, numColumns)
	values := make([][]any, len(q.Values))

	for i, m := range q.Values {
		if len(q.Values[i]) != numColumns {
			return dialects.RawQuery{}, ErrInsertMismatchedValueKeys
		}
		values[i] = make([]any, 0, numColumns)
		for k, v := range m {
			if i == 0 {
				columns = append(columns, k)
			}
			values[i] = append(values[i], v)
		}
	}

	b := newRawQueryBuilder().
		AddString("INSERT INTO").
		AddString(g.core.Identifier(q.Table)).
		AddString("(" + strings.Join(slices.Map(columns, g.core.Identifier), ", ") + ")").
		AddString("VALUES")

	for i, v := range values {
		if i > 0 {
			b.AddStringNoSpace(",")
		}
		b.Add(group(mapJoinRawQueries(v, ", ", g.EncodeLiteral)))
	}
	return b.Build()
}
