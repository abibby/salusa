package dialects

type MultiQueryBuilder interface {
	MultiQuery() []MultiQuery
}

type MultiQuery struct {
	alterTableQuery  *AlterTableQuery
	createTableQuery *CreateTableQuery
	deleteQuery      *DeleteQuery
	dropTableQuery   *DropTableQuery
	insertQuery      *InsertQuery
	selectQuery      *SelectQuery
	updateQuery      *UpdateQuery
}

type MultiQueryBuilderImpl struct {
	queries []MultiQuery
}

func NewMultiQueryBuilder() *MultiQueryBuilderImpl {
	return &MultiQueryBuilderImpl{
		queries: []MultiQuery{},
	}
}

func (b *MultiQueryBuilderImpl) AddAlterTableQuery(q *AlterTableQuery) *MultiQueryBuilderImpl {
	b.queries = append(b.queries, MultiQuery{alterTableQuery: q})
	return b
}
func (b *MultiQueryBuilderImpl) AddCreateTableQuery(q *CreateTableQuery) *MultiQueryBuilderImpl {
	b.queries = append(b.queries, MultiQuery{createTableQuery: q})
	return b
}
func (b *MultiQueryBuilderImpl) AddDeleteQuery(q *DeleteQuery) *MultiQueryBuilderImpl {
	b.queries = append(b.queries, MultiQuery{deleteQuery: q})
	return b
}
func (b *MultiQueryBuilderImpl) AddDropTableQuery(q *DropTableQuery) *MultiQueryBuilderImpl {
	b.queries = append(b.queries, MultiQuery{dropTableQuery: q})
	return b
}
func (b *MultiQueryBuilderImpl) AddInsertQuery(q *InsertQuery) *MultiQueryBuilderImpl {
	b.queries = append(b.queries, MultiQuery{insertQuery: q})
	return b
}
func (b *MultiQueryBuilderImpl) AddSelectQuery(q *SelectQuery) *MultiQueryBuilderImpl {
	b.queries = append(b.queries, MultiQuery{selectQuery: q})
	return b
}
func (b *MultiQueryBuilderImpl) AddUpdateQuery(q *UpdateQuery) *MultiQueryBuilderImpl {
	b.queries = append(b.queries, MultiQuery{updateQuery: q})
	return b
}

func (b *MultiQueryBuilderImpl) MultiQuery(q *UpdateQuery) []MultiQuery {
	return b.queries
}

func EncodeMultiQuery(d Dialect, queries []MultiQuery) (RawQuery, error) {
	rawQueries := make([]RawQuery, len(queries))
	var err error
	for i, q := range queries {
		rawQueries[i], err = encodeMultiQuery(d, q)
		if err != nil {
			return RawQuery{}, err
		}
	}
	return JoinQueries(rawQueries), nil
}
func encodeMultiQuery(d Dialect, q MultiQuery) (RawQuery, error) {
	if q.alterTableQuery != nil {
		return d.EncodeAlterTableQuery(q.alterTableQuery)
	} else if q.createTableQuery != nil {
		return d.EncodeCreateTableQuery(q.createTableQuery)
	} else if q.deleteQuery != nil {
		return d.EncodeDeleteQuery(q.deleteQuery)
	} else if q.dropTableQuery != nil {
		return d.EncodeDropTableQuery(q.dropTableQuery)
	} else if q.insertQuery != nil {
		return d.EncodeInsertQuery(q.insertQuery)
	} else if q.selectQuery != nil {
		return d.EncodeSelectQuery(q.selectQuery)
	} else if q.updateQuery != nil {
		return d.EncodeUpdateQuery(q.updateQuery)
	}
	return RawQuery{}, nil
}
