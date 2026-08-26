package builder

import (
	"abibby.com/salusa/database/dialects"
)

// Select sets the columns to be selected.
func (b *Builder) Select(columns ...string) *Builder {
	identifiers := make([]dialects.Column, len(columns))
	for i, c := range columns {
		identifiers[i] = dialects.Column{Column: c}
	}
	b.query.Select.Columns = identifiers
	return b
}

// AddSelect adds new columns to be selected.
func (b *Builder) AddSelect(columns ...string) *Builder {
	for _, c := range columns {
		b.query.Select.Columns = append(b.query.Select.Columns, dialects.Column{Column: c})
	}
	return b
}

// SelectSubquery sets a subquery to be selected.
func (b *Builder) SelectSubquery(sb dialects.QueryBuilder, as string) *Builder {
	return b.Select().AddSelectSubquery(sb, as)
}

// AddSelectSubquery adds a subquery to be selected.
func (b *Builder) AddSelectSubquery(sb dialects.QueryBuilder, as string) *Builder {
	b.query.Select.Columns = append(b.query.Select.Columns, dialects.Column{
		SubQuery: sb,
		As:       as,
	})
	return b
}

// SelectFunction sets a column to be selected with a function applied.
func (b *Builder) SelectFunction(function, column string) *Builder {
	return b.Select().AddSelectFunction(function, column)
}

// SelectFunction adds a column to be selected with a function applied.
func (b *Builder) AddSelectFunction(function, column string) *Builder {
	b.query.Select.Columns = append(b.query.Select.Columns, dialects.Column{
		Function: &dialects.FunctionCall{
			Name:      function,
			Arguments: column,
		},
	})

	return b
}

// Distinct forces the query to only return distinct results.
func (b *Builder) Distinct() *Builder {
	b.query.Select.Distinct = true
	return b
}
