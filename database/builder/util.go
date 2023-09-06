package builder

import (
	"github.com/abibby/salusa/database/dialects"
	"github.com/davecgh/go-spew/spew"
)

func (b *SubBuilder) Dump() *SubBuilder {
	spew.Dump(b.ToSQL(dialects.New()))
	return b
}
