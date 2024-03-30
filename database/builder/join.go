package builder

import (
	"github.com/abibby/salusa/database/dialects"
	"github.com/abibby/salusa/internal/helpers"
)

type join struct {
	direction  string
	table      helpers.SQLStringer
	conditions *Conditions
}

type joins []*join

func (j joins) Clone() joins {
	return cloneSlice(j)
}

// Join adds a join clause to the query.
func (j joins) Join(table, localColumn, operator, foreignColumn string) joins {
	return j.join("", table, localColumn, operator, foreignColumn)
}

// LeftJoin adds a left join clause to the query.
func (j joins) LeftJoin(table, localColumn, operator, foreignColumn string) joins {
	return j.join("LEFT", table, localColumn, operator, foreignColumn)
}

// RightJoin adds a right join clause to the query.
func (j joins) RightJoin(table, localColumn, operator, foreignColumn string) joins {
	return j.join("RIGHT", table, localColumn, operator, foreignColumn)
}

// InnerJoin adds an inner join clause to the query.
func (j joins) InnerJoin(table, localColumn, operator, foreignColumn string) joins {
	return j.join("INNER", table, localColumn, operator, foreignColumn)
}

// CrossJoin adds a cross join clause to the query.
func (j joins) CrossJoin(table, localColumn, operator, foreignColumn string) joins {
	return j.join("CROSS", table, localColumn, operator, foreignColumn)
}
func (j joins) join(direction, table, localColumn, operator, foreignColumn string) joins {
	return j.joinOn(direction, table, func(q *Conditions) {
		q.WhereColumn(localColumn, operator, foreignColumn)
	})
}

// JoinOn adds a join clause to the query with a complex on statement.
func (j joins) JoinOn(table string, cb func(q *Conditions)) joins {
	return j.joinOn("", table, cb)
}

// LeftJoinOn adds a left join clause to the query with a complex on statement.
func (j joins) LeftJoinOn(table string, cb func(q *Conditions)) joins {
	return j.joinOn("LEFT", table, cb)
}

// RightJoinOn adds a right join clause to the query with a complex on statement.
func (j joins) RightJoinOn(table string, cb func(q *Conditions)) joins {
	return j.joinOn("RIGHT", table, cb)
}

// InnerJoinOn adds an inner join clause to the query with a complex on statement.
func (j joins) InnerJoinOn(table string, cb func(q *Conditions)) joins {
	return j.joinOn("INNER", table, cb)
}

// CrossJoinOn adds a cross join clause to the query with a complex on statement.
func (j joins) CrossJoinOn(table string, cb func(q *Conditions)) joins {
	return j.joinOn("CROSS", table, cb)
}
func (j joins) joinOn(direction string, table string, cb func(q *Conditions)) joins {
	c := newConditions().withPrefix("ON")
	cb(c)
	return append(j, &join{
		direction:  direction,
		table:      helpers.Identifier(table),
		conditions: c,
	})
}
func (j *join) SQLString(d dialects.Dialect) (string, []any, error) {
	r := helpers.Result().
		AddString(j.direction).
		AddString("JOIN").
		Add(j.table).
		Add(j.conditions)

	return r.SQLString(d)
}
func (j joins) SQLString(d dialects.Dialect) (string, []any, error) {
	return helpers.Join(j, " ").SQLString(d)
}
