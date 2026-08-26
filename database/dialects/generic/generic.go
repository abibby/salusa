package generic

import (
	"errors"

	"abibby.com/salusa/database/dialects"
)

var ErrUnkownExprType = errors.New("unknown dialects.Expr type")

type Core interface {
	Identifier(string) string
	DataType(dialects.DataType) string
	CurrentTime() string
	AutoIncrement() string
	Escape(v any) string
	Binding() string
	Features() dialects.Features
}
type Generic struct {
	core Core
}

var _ dialects.Dialect = (*Generic)(nil)

func New(c Core) *Generic {
	return &Generic{core: c}
}

// Features implements [dialects.Dialect].
func (g *Generic) Features() dialects.Features {
	return g.core.Features()
}
