package dialects

type UopdateQuery struct {
	Table  string
	Values map[string]any
	Wheres []Condition
}
