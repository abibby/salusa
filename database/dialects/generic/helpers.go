package generic

import "abibby.com/salusa/database/dialects"

func group(r dialects.RawQuery, err error) (dialects.RawQuery, error) {
	r.SQL = "(" + r.SQL + ")"
	return r, err
}
