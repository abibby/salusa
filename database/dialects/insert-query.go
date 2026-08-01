package dialects

type InsertQuery struct {
	Table  string
	Values []map[string]any
}
