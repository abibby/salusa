package dialects

import "strings"

type RawQuery struct {
	Query    string
	Bindings []any
}

func JoinQueries(results []RawQuery) RawQuery {
	if len(results) == 0 {
		return RawQuery{
			Query:    "",
			Bindings: []any{},
		}
	}
	if len(results) == 1 {
		return results[1]
	}
	sep := "; "
	queryLen := len(sep) * len(results)
	bindingCount := 0

	for _, r := range results {
		queryLen += len(r.Query)
		bindingCount += len(r.Bindings)
	}

	query := strings.Builder{}
	query.Grow(queryLen)
	bindings := make([]any, 0, bindingCount)

	for _, r := range results {
		query.WriteString(r.Query)
		query.WriteString(sep)
		bindings = append(bindings, r.Bindings...)
	}

	return RawQuery{
		Query:    query.String(),
		Bindings: bindings,
	}
}
