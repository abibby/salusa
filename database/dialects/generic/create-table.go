package generic

import (
	"github.com/abibby/salusa/database/dialects"
)

func (g *Generic) EncodeCreateTableQuery(q *dialects.CreateTableQuery) (dialects.SQLResult, error) {
	b := resultBuilder().AddString("CREATE TABLE")
	if q.IfNotExists {
		b.AddString("IF NOT EXISTS")
	}
	b.AddString(g.core.Identifier(q.Table))

	var err error
	columns := make([]dialects.SQLResult, len(q.Columns))
	for i, c := range q.Columns {
		columns[i], err = g.EncodeColumnDefinition(&c)
		if err != nil {
			return dialects.SQLResult{}, err
		}
		columns[i].Query = "\n\t" + columns[i].Query
		if i == len(q.Columns)-1 {
			columns[i].Query += "\n"
		}

	}
	b.Add(group(joinResults(columns, ","), nil))

	return b.Build()
}

func (g *Generic) EncodeColumnDefinition(c *dialects.ColumnDefinition) (dialects.SQLResult, error) {
	r := resultBuilder()
	r.AddString(g.core.Identifier(c.Name))
	r.AddString(g.core.DataType(c.Datatype))

	if c.AutoIncrement {
		r.AddString("PRIMARY KEY " + g.core.AutoIncrement())
	} else if c.Primary {
		r.AddString("PRIMARY KEY")
	}
	if !c.Nullable {
		r.AddString("NOT NULL")
	}
	if c.Unique {
		r.AddString("UNIQUE")
	}

	if c.DefaultValue != nil {
		r.AddString("DEFAULT").
			AddString(g.core.Escape(c.DefaultValue))
	} else if c.DefaultCurrentTime {
		r.AddString("DEFAULT").
			AddString(g.core.CurrentTime())
	}
	return r.Build()
}
