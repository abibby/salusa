package dialects

type InsertQueryBuilder interface {
	InsertQuery() *InsertQuery
}
type InsertQuery struct {
	Table  string
	Values []map[string]any
}
