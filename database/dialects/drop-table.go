package dialects

type DropTableQuery struct {
	Table    string
	IfExists bool
}
