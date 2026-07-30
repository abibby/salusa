package dialects

type QueryBuilder interface {
	Query() *Query
}

type Query struct {
	Select Select
	From   string
	// Joins    []Join
	Wheres  []Condition
	Havings []Condition
	// GroupBys []string
	// OrderBys orderBys
	// Limit    int
	// Offset   int
}

type Select struct {
	Distinct bool
	Columns  []Column
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
