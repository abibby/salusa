package dialects

type UpdateQuery struct {
	Table  string
	Values map[string]any
	Wheres []Condition
}
