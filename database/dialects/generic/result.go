package generic

import (
	"strings"

	"github.com/abibby/salusa/database/dialects"
)

type sqlResultBuilder struct {
	results []dialects.SQLResult
	err     error
}

func resultBuilder() *sqlResultBuilder {
	return &sqlResultBuilder{
		results: []dialects.SQLResult{},
	}
}

func (b *sqlResultBuilder) Add(r dialects.SQLResult, err error) *sqlResultBuilder {
	if b.err != nil || r.Query == "" {
		return b
	}
	b.results = append(b.results, r)
	b.err = err
	return b
}

func (b *sqlResultBuilder) AddString(s string) *sqlResultBuilder {
	return b.Add(dialects.SQLResult{
		Query: s,
	}, nil)
}

func (b *sqlResultBuilder) Build() (dialects.SQLResult, error) {
	if b.err != nil {
		return dialects.SQLResult{}, b.err
	}
	return joinResults(b.results, " "), nil
}

func joinResults(results []dialects.SQLResult, sep string) dialects.SQLResult {
	if len(results) == 0 {
		return dialects.SQLResult{
			Query:    "",
			Bindings: []any{},
		}
	}
	if len(results) == 1 {
		return results[0]
	}

	queryLen := len(sep) * (len(results) - 1)
	bindingCount := 0

	for _, r := range results {
		queryLen += len(r.Query)
		bindingCount += len(r.Bindings)
	}

	query := strings.Builder{}
	query.Grow(queryLen)
	bindings := make([]any, 0, bindingCount)

	query.WriteString(results[0].Query)
	bindings = append(bindings, results[0].Bindings...)

	for _, r := range results[1:] {
		query.WriteString(sep)
		query.WriteString(r.Query)
		bindings = append(bindings, r.Bindings...)
	}

	return dialects.SQLResult{
		Query:    query.String(),
		Bindings: bindings,
	}
}

func mapJoinResults[T any](arr []T, sep string, fn func(v T) (dialects.SQLResult, error)) (dialects.SQLResult, error) {
	results := make([]dialects.SQLResult, len(arr))
	var err error
	for i, v := range arr {
		results[i], err = fn(v)
		if err != nil {
			return dialects.SQLResult{}, err
		}
	}
	return joinResults(results, sep), nil
}
