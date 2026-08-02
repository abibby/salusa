package dialects

type QueryBuilder interface {
	Query() *SelectQuery
}

type SelectQuery struct {
	Select   Select
	From     string
	Joins    []Join
	Wheres   []Condition
	Havings  []Condition
	GroupBys []string
	OrderBys []OrderColumn
	Limit    Limit
}

func NewSelectQuery() SelectQuery {
	return SelectQuery{
		Select:   NewSelect(),
		Joins:    []Join{},
		Wheres:   []Condition{},
		GroupBys: []string{},
		Havings:  []Condition{},
		OrderBys: []OrderColumn{},
	}
}

func (q *SelectQuery) Clone() *SelectQuery {
	// TODO implement clone
	return q
}

type OrderColumn struct {
	Column     string
	Descending bool
}

type Select struct {
	Distinct bool
	Columns  []Column
}

func NewSelect() Select {
	return Select{
		Columns: []Column{},
	}
}

type Join struct {
	Direction  string
	Table      string
	Conditions []Condition
}

type Condition struct {
	Column   Column
	Operator string
	Value    any
	Or       bool
}

type Column struct {
	Column   string
	Function *FunctionCall
	SubQuery QueryBuilder

	As string
}

type FunctionCall struct {
	Name      string
	Arguments string
}

type Limit struct {
	Limit  int
	Offset int
}

type RawString string

func Raw(sql string, bindings ...any) RawQuery {
	return RawQuery{
		SQL:      sql,
		Bindings: bindings,
	}
}
