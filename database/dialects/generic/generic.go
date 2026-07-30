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
	// Identifier    func(string) string
	// DataType      func(dialects.DataType) string
	// CurrentTime   func() string
	// AutoIncrement func() string
	// Escape        func(v any) string
	// Binding       func() string
}

func New(c Core) Generic {
	return Generic{core: c}
}
