package dialects

type Condition struct {
	Column   Column
	Operator string
	Value    any
	Or       bool
}
