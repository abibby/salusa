package builder

import (
	"github.com/abibby/salusa/database"
	"github.com/abibby/salusa/database/dialects"
	"github.com/abibby/salusa/di"
	"github.com/abibby/salusa/salusaconfig"
	"github.com/davecgh/go-spew/spew"
)

func (b *Builder) Dump() *Builder {
	cfg, err := di.Resolve[salusaconfig.Config](b.ctx)
	if err != nil {
		return b
	}
	dbcfg, ok := cfg.(database.DBConfiger)
	if !ok {
		return b
	}
	d, err := dialects.New(dbcfg.DBConfig().DriverName())
	if err != nil {
		return b
	}
	spew.Dump(d.EncodeSelectQuery(b.Query()))
	return b
}
