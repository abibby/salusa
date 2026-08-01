package dialects

type DeleteQuery struct {
	Table  string
	Wheres []Condition
}
