package generic

import "github.com/abibby/salusa/database/dialects"

func (g *Generic) EncodeLimit(l *dialects.Limit) (dialects.RawQuery, error) {
	if l.Limit == 0 && l.Offset == 0 {
		return dialects.RawQuery{}, nil
	}
	b := newRawQueryBuilder()
	if l.Limit != 0 {
		b.AddString("LIMIT").Add(g.EncodeLiteral(l.Limit))
	}
	if l.Offset != 0 {
		b.AddString("OFFSET").Add(g.EncodeLiteral(l.Offset))
	}
	return b.Build()
}
