package builder

import (
	"github.com/abibby/salusa/database/dialects"
	"github.com/davecgh/go-spew/spew"
)

func (b *Builder) Dump() *Builder {
	spew.Dump(dialects.New().EncodeSelectQuery(b.Query()))
	return b
}
