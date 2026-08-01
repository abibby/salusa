package generic

import "github.com/abibby/salusa/database/dialects"

func group(r dialects.RawQuery, err error) (dialects.RawQuery, error) {
	r.Query = "(" + r.Query + ")"
	return r, err
}
