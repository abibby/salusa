package builder

import (
	"context"
	"fmt"
	"reflect"

	"github.com/abibby/salusa/database/dialects"
	"github.com/abibby/salusa/internal/helpers"
)

type where struct {
	Column   helpers.ToSQLer
	Operator string
	Value    any
	Or       bool
}

type Conditions struct {
	parent any
	prefix string
	list   []*where
	ctx    context.Context
}

func newConditions() *Conditions {
	return &Conditions{
		list: []*where{},
		ctx:  context.Background(),
	}
}

func (c *Conditions) withPrefix(prefix string) *Conditions {
	c.prefix = prefix
	return c
}
func (c *Conditions) withParent(parent any) *Conditions {
	c.parent = parent
	return c
}
func (c *Conditions) Clone() *Conditions {
	return &Conditions{
		parent: c.parent,
		prefix: c.prefix,
		list:   cloneSlice(c.list),
		ctx:    c.ctx,
	}
}
func (c *Conditions) ToSQL(d dialects.Dialect) (string, []any, error) {
	if len(c.list) == 0 {
		return "", nil, nil
	}

	r := helpers.Result()
	if c.prefix != "" {
		r.AddString(c.prefix)
	}
	for i, c := range c.list {
		if i != 0 {
			if c.Or {
				r.AddString("OR")
			} else {
				r.AddString("AND")
			}
		}
		if c.Column != nil {
			r.Add(c.Column)

			if c.Operator == "" {
				return "", nil, fmt.Errorf("the operator must be set when the column is set")
			}
		}

		if c.Value == nil {
			switch c.Operator {
			case "=":
				r.AddString("IS NULL")
			case "!=":
				r.AddString("IS NOT NULL")
			default:
				return "", nil, fmt.Errorf("wheres checking nil only support = and !=")
			}
		} else {
			if c.Operator != "" {
				r.AddString(c.Operator)
			}
			if sb, ok := c.Value.(QueryBuilder); ok {
				r.Add(helpers.Group(sb))
			} else if sb, ok := c.Value.(*Conditions); ok {
				r.Add(helpers.Group(sb))
			} else if sb, ok := c.Value.(helpers.ToSQLer); ok {
				r.Add(sb)
			} else {
				r.Add(helpers.Literal(c.Value))
			}
		}
	}

	return r.ToSQL(d)
}

// Where adds a basic where clause to the query.
func (c *Conditions) Where(column, operator string, value any) *Conditions {
	return c.where(column, operator, value, false)
}

// OrWhere adds an or where clause to the query
func (c *Conditions) OrWhere(column, operator string, value any) *Conditions {
	return c.where(column, operator, value, true)
}

// WhereColumn adds a where clause to the query comparing two columns.
func (c *Conditions) WhereColumn(column, operator string, valueColumn string) *Conditions {
	return c.where(column, operator, helpers.Identifier(valueColumn), false)
}

// OrWhereColumn adds an or where clause to the query comparing two columns.
func (c *Conditions) OrWhereColumn(column, operator string, valueColumn string) *Conditions {
	return c.where(column, operator, helpers.Identifier(valueColumn), true)
}

// WhereIn adds a where in clause to the query.
func (c *Conditions) WhereIn(column string, values []any) *Conditions {
	return c.whereIn(column, values, false)
}

// OrWhereIn adds an or where in clause to the query.
func (c *Conditions) OrWhereIn(column string, values []any) *Conditions {
	return c.whereIn(column, values, true)
}

func (c *Conditions) whereIn(column string, values []any, or bool) *Conditions {
	return c.where(column, "in", helpers.Group(helpers.Join(helpers.LiteralList(values), ", ")), or)
}

// WhereExists add an exists clause to the query.
func (c *Conditions) WhereExists(query QueryBuilder) *Conditions {
	return c.whereExists(query, false)
}

// WhereExists add an exists clause to the query.
func (c *Conditions) OrWhereExists(query QueryBuilder) *Conditions {
	return c.whereExists(query, true)
}

func (c *Conditions) whereExists(query QueryBuilder, or bool) *Conditions {
	return c.addWhere(&where{
		Value: helpers.Join([]helpers.ToSQLer{
			helpers.Raw("EXISTS"),
			helpers.Group(query),
		}, " "),
		Or: or,
	})
}

// WhereSubquery adds a where clause to the query comparing a column and a subquery.
func (c *Conditions) WhereSubquery(subquery QueryBuilder, operator string, value any) *Conditions {
	return c.whereSubquery(subquery, operator, value, false)
}

// OrWhereSubquery adds an or where clause to the query comparing a column and a subquery.
func (c *Conditions) OrWhereSubquery(subquery QueryBuilder, operator string, value any) *Conditions {
	return c.whereSubquery(subquery, operator, value, true)
}

func (c *Conditions) whereSubquery(subquery QueryBuilder, operator string, value any, or bool) *Conditions {
	return c.addWhere(&where{
		Column:   helpers.Group(subquery),
		Operator: operator,
		Value:    value,
		Or:       or,
	})
}

func (c *Conditions) where(column, operator string, value any, or bool) *Conditions {
	return c.addWhere(&where{
		Column:   helpers.Identifier(column),
		Operator: operator,
		Value:    value,
		Or:       or,
	})
}

// WhereHas adds a relationship exists condition to the query with where clauses.
func (c *Conditions) WhereHas(relation string, cb func(q *SubBuilder) *SubBuilder) *Conditions {
	return c.whereHas(relation, cb, false)
}

// OrWhereHas adds a relationship exists condition to the query with where clauses and an or.
func (c *Conditions) OrWhereHas(relation string, cb func(q *SubBuilder) *SubBuilder) *Conditions {
	return c.whereHas(relation, cb, true)
}
func (c *Conditions) whereHas(relation string, cb func(q *SubBuilder) *SubBuilder, or bool) *Conditions {
	r, ok := getRelation(reflect.ValueOf(c.parent), relation)
	if !ok {
		return c
	}

	return c.whereExists(cb(r.Subquery().WithContext(c.ctx)), or)
}

// WhereRaw adds a raw where clause to the query.
func (c *Conditions) WhereRaw(rawSql string, bindings ...any) *Conditions {
	return c.whereRaw(rawSql, bindings, false)
}

// OrWhereRaw adds a raw or where clause to the query.
func (c *Conditions) OrWhereRaw(rawSql string, bindings ...any) *Conditions {
	return c.whereRaw(rawSql, bindings, true)
}
func (c *Conditions) whereRaw(rawSql string, bindings []any, or bool) *Conditions {
	return c.addWhere(&where{
		Value: helpers.Raw(rawSql, bindings...),
		Or:    or,
	})
}

// And adds a group of conditions to the query
func (c *Conditions) And(cb func(q *Conditions)) *Conditions {
	return c.group(cb, false)
}

// Or adds a group of conditions to the query with an or
func (c *Conditions) Or(cb func(q *Conditions)) *Conditions {
	return c.group(cb, true)
}

func (c *Conditions) group(cb func(wl *Conditions), or bool) *Conditions {
	wl := newConditions().withParent(c.parent)
	wl.ctx = c.ctx
	cb(wl)
	return c.addWhere(&where{Value: wl, Or: or})
}

func (c *Conditions) addWhere(wh *where) *Conditions {
	c.list = append(c.list, wh)
	return c
}
