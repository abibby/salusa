package dialects

type UpdateQueryBuilder interface {
	UpdateQuery() *UpdateQuery
}
type UpdateQuery struct {
	Table  string
	Values map[string]any
	Wheres []Condition
}
