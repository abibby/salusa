package dialects

type Expr interface {
}

type Literal struct {
	Value any
}
type FunctionCall struct {
	Name       string
	Parameters []Expr
}
