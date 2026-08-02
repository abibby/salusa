package dialects

type DropTableQueryBuilder interface {
	DropTableQuery() *DropTableQuery
}
type DropTableQuery struct {
	Table    string
	IfExists bool
}
