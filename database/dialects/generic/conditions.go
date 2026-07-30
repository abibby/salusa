package generic

import (
	"fmt"

	"github.com/abibby/salusa/database/dialects"
)

func (g *Generic) EncodeWheres(c []dialects.Condition) (dialects.SQLResult, error) {
	return g.encodeConditionsPrefix("WHERE", c)
}
func (g *Generic) EncodeHavings(c []dialects.Condition) (dialects.SQLResult, error) {
	return g.encodeConditionsPrefix("HAVING", c)
}

func (g *Generic) encodeConditionsPrefix(prefix string, c []dialects.Condition) (dialects.SQLResult, error) {
	if len(c) == 0 {
		return dialects.SQLResult{}, nil
	}
	return ResultBuilder().AddString(prefix).Add(g.EncodeConditions(c)).Build()
}

func (g *Generic) EncodeConditions(c []dialects.Condition) (dialects.SQLResult, error) {
	b := ResultBuilder()
	for i, c := range c {
		if i != 0 {
			if c.Or {
				b.AddString("OR")
			} else {
				b.AddString("AND")
			}
		}
		if (c.Column != dialects.Column{}) {
			b.Add(g.EncodeColumn(&c.Column))

			if c.Operator == "" {
				return dialects.SQLResult{}, fmt.Errorf("the operator must be set when the column is set")
			}
		}

		if c.Value == nil {
			switch c.Operator {
			case "=":
				b.AddString("IS NULL")
			case "!=":
				b.AddString("IS NOT NULL")
			default:
				return dialects.SQLResult{}, fmt.Errorf("wheres checking nil only support = and !=")
			}
		} else {
			if c.Operator != "" {
				b.AddString(c.Operator)
			}
			if sb, ok := c.Value.(dialects.QueryBuilder); ok {
				b.Add(Group(g.EncodeQuery(sb.Query())))
			} else if sb, ok := c.Value.([]dialects.Condition); ok {
				b.Add(Group(g.EncodeConditions(sb)))
			} else {
				b.Add(g.EncodeLiteral(c.Value))
			}
		}
	}

	return b.Build()
}
