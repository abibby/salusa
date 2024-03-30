package builder

import (
	"github.com/abibby/salusa/database/dialects"
	"github.com/davecgh/go-spew/spew"
)

func (b *Builder) Dump() *Builder {
	spew.Dump(b.ToSQL(dialects.New()))
	return b
}
