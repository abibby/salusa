package dialects

type Index struct {
	Name    string
	Table   string
	Columns []string
	Unique  bool
}
