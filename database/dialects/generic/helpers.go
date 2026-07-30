package generic

import "github.com/abibby/salusa/database/dialects"

func Group(r dialects.SQLResult, err error) (dialects.SQLResult, error) {
	r.Query = "(" + r.Query + ")"
	return r, err
}
