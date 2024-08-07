package builder

import (
	"fmt"
	"strings"

	"github.com/abibby/salusa/database/dialects"
	"github.com/abibby/salusa/internal/helpers"
)

type selects struct {
	distinct bool
	list     []helpers.SQLStringer
}

func NewSelects() *selects {
	return &selects{
		list: []helpers.SQLStringer{},
	}
}

func (s *selects) Clone() *selects {
	return &selects{
		distinct: s.distinct,
		list:     cloneSlice(s.list),
	}
}
func (s *selects) SQLString(d dialects.Dialect) (string, []any, error) {
	if len(s.list) == 0 {
		return "", nil, nil
	}
	r := helpers.Result()
	r.AddString("SELECT")
	if s.distinct {
		r.AddString("DISTINCT")
	}
	r.Add(helpers.Join(s.list, ", "))
	return r.SQLString(d)
}

// Select sets the columns to be selected.
func (s *selects) Select(columns ...string) *selects {
	identifiers := make([]helpers.SQLStringer, len(columns))
	for i, c := range columns {
		identifiers[i] = columnIdentifier(c)
	}
	s.list = identifiers
	return s
}

// AddSelect adds new columns to be selected.
func (s *selects) AddSelect(columns ...string) *selects {
	for _, c := range columns {
		s.list = append(s.list, columnIdentifier(c))
	}
	return s
}

// SelectSubquery sets a subquery to be selected.
func (s *selects) SelectSubquery(sb QueryBuilder, as string) *selects {
	return s.Select().AddSelectSubquery(sb, as)
}

// AddSelectSubquery adds a subquery to be selected.
func (s *selects) AddSelectSubquery(sb QueryBuilder, as string) *selects {
	s.list = append(s.list, helpers.Concat(
		helpers.Group(sb),
		helpers.Raw(" as "),
		helpers.Identifier(as),
	))

	return s
}

// SelectFunction sets a column to be selected with a function applied.
func (s *selects) SelectFunction(function, column string) *selects {
	return s.Select().AddSelectFunction(function, column)
}

// SelectFunction adds a column to be selected with a function applied.
func (s *selects) AddSelectFunction(function, column string) *selects {
	s.list = append(s.list, helpers.SQLStringFunc(func(d dialects.Dialect) (string, []any, error) {
		var c helpers.SQLStringer
		if column == "*" {
			c = helpers.Raw("*")
		} else {
			c = helpers.Identifier(column)
		}
		q, bindings, err := c.SQLString(d)
		if err != nil {
			return "", nil, err
		}
		return fmt.Sprintf("%s(%s)", function, q), bindings, nil
	}))

	return s
}

// Distinct forces the query to only return distinct results.
func (s *selects) Distinct() *selects {
	s.distinct = true
	return s
}

func columnIdentifier(column string) helpers.SQLStringer {
	var identifier helpers.SQLStringer
	parts := strings.SplitN(column, " as ", 2)
	c := parts[0]
	if c == "*" {
		identifier = helpers.Raw("*")
	} else if strings.HasSuffix(c, ".*") {
		identifier = helpers.Concat(helpers.Identifier(c[:len(c)-2]), helpers.Raw(".*"))
	} else {
		identifier = helpers.Identifier(c)
	}
	if len(parts) > 1 {
		identifier = helpers.Concat(
			identifier,
			helpers.Raw(" as "),
			helpers.Identifier(parts[1]),
		)
	}
	return identifier
}
