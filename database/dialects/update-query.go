package dialects

type UpdateQueryBuilder interface {
	UpdateQuery() *UpdateQuery
}
type UpdateQuery struct {
	Table     string
	Values    map[string]any
	Returning []Column
	Wheres    []Condition
}
