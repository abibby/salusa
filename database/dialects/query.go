package dialects

type Query struct {
	Selects Selects
	From    string
	// Joins    []Join
	// Wheres   Conditions
	// Havings  Conditions
	// GroupBys []string
	// OrderBys orderBys
	// Limit    int
	// Offset   int
}

type Selects struct {
	Distinct bool
	Columns  []Column
}

type Column struct {
	Expr Expr
	As   string
}
