package dialects

type QueryBuilder interface {
	Query() *Query
}

type Query struct {
	Select   Select
	From     string
	Joins    []Join
	Wheres   []Condition
	Havings  []Condition
	GroupBys []string
	OrderBys []string
	Limit    Limit
}

type Select struct {
	Distinct bool
	Columns  []Column
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

type Raw string
