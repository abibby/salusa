package generic

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/abibby/salusa/database/dialects"
)

var ErrUnkownExprType = errors.New("unknown dialects.Expr type")

type Generic struct {
	Identifier    func(string) string
	DataType      func(dt dialects.DataType) string
	CurrentTime   func() string
	AutoIncrement func() string
	Escape        func(v any) string
	Binding       func() string
}

func (g *Generic) EncodeExpr(e dialects.Expr) (dialects.SQLResult, error) {
	switch e := e.(type) {
	case *dialects.Query:
		return g.EncodeQuery(e)
	case *dialects.Literal:
		return g.EncodeLiteral(e)
	}
	return dialects.SQLResult{}, fmt.Errorf("%w: %s", ErrUnkownExprType, reflect.TypeOf(e).Name())
}

func (g *Generic) EncodeLiteral(l *dialects.Literal) (dialects.SQLResult, error) {
	return dialects.SQLResult{
		Query:    g.Binding(),
		Bindings: []any{l.Value},
	}, nil
}

func (g *Generic) EncodeFunctionCall(fc *dialects.FunctionCall) (dialects.SQLResult, error) {
	params := make([]dialects.SQLResult, len(fc.Parameters))
	var err error

	for i, p := range fc.Parameters {
		params[i], err = g.EncodeExpr(p)
		if err != nil {
			return dialects.SQLResult{}, err
		}
	}

	return dialects.ResultBuilder().
		AddString(fc.Name+"(").
		Add(dialects.JoinResults(params, ", "), nil).
		AddString(")").
		Build()
}

func (g *Generic) EncodeQuery(q *dialects.Query) (dialects.SQLResult, error) {
	// return helpers.Result().
	// 	Add(b.selects).
	// 	Add(b.from).
	// 	Add(b.joins).
	// 	Add(b.wheres).
	// 	Add(b.groupBys).
	// 	Add(b.havings).
	// 	Add(b.orderBys).
	// 	Add(b.limit).
	// 	SQLString(d)
	// JoinResults()

	return dialects.ResultBuilder().
		Add(g.EncodeSelects(&q.Selects)).
		Add(g.EncodeFrom(q.From)).
		Build()
}

func (g *Generic) EncodeSelects(s *dialects.Selects) (dialects.SQLResult, error) {
	if len(s.Columns) == 0 {
		return dialects.SQLResult{}, nil
	}

	b := dialects.ResultBuilder()
	b.AddString("SELECT")
	if s.Distinct {
		b.AddString("DISTINCT")
	}

	for _, c := range s.Columns {
		b.Add(g.EncodeColumn(&c))
	}

	return b.Build()
}
func (g *Generic) EncodeFrom(from string) (dialects.SQLResult, error) {
	if from == "" {
		return dialects.SQLResult{}, nil
	}
	return dialects.SQLResult{
		Query: "FROM " + from,
	}, nil
}

func (g *Generic) EncodeColumn(e *dialects.Column) (dialects.SQLResult, error) {
	return dialects.ResultBuilder().
		Add(g.EncodeExpr(&e.Expr)).
		AddString(e.As).
		Build()
}
