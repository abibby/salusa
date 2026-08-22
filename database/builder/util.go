package builder

import (
	"github.com/abibby/salusa/database/dialects"
	"github.com/abibby/salusa/di"
	"github.com/davecgh/go-spew/spew"
	"github.com/jmoiron/sqlx"
)

func (b *Builder) Dump() *Builder {
	db, err := di.Resolve[*sqlx.DB](b.ctx)
	if err != nil {
		return b
	}
	d, err := dialects.New(db.DriverName())
	if err != nil {
		return b
	}
	spew.Dump(d.EncodeSelectQuery(b.Query()))
	return b
}
