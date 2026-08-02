package generic

import (
	"errors"

	"github.com/abibby/salusa/database/dialects"
)

var ErrUnkownExprType = errors.New("unknown dialects.Expr type")

type Core interface {
	Identifier(string) string
	DataType(dialects.DataType) string
	CurrentTime() string
	AutoIncrement() string
	Escape(v any) string
	Binding() string
}
type Generic struct {
	core Core
}

func New(c Core) *Generic {
	return &Generic{core: c}
}
