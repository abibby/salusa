package dialects

type DeleteQueryBuilder interface {
	DeleteQuery() *DeleteQuery
}
type DeleteQuery struct {
	Table  string
	Wheres []Condition
}
