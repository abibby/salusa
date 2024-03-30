package builder

import (
	"github.com/abibby/salusa/database/dialects"
	"github.com/davecgh/go-spew/spew"
)

func (b *Builder) Dump() *Builder {
	spew.Dump(b.SQLString(dialects.New()))
	return b
}
