package generateclient

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/abibby/salusa/router"
)

func GenerateClient(r *router.Router, w io.Writer) error {
	for _, route := range r.Routes() {
		err := routeTs(route, w)
		if err != nil {
			return err
		}
	}
	return nil
}

var funcTpl = `export async function %s(): Promise<unknown> {
    const response = await fetcher(%s)
    if (response.status < 200 || response.status >= 300) {
        throw new Error("invalid status")
    }
    return await response.json()
}`

func jsonString(s string) string {
	b, err := json.Marshal(s)
	if err != nil {
		panic(err)
	}
	return string(b)
}
func routeTs(r *router.Route, w io.Writer) error {
	if ht, ok := r.GetHandler().(HandlerTypes); ok {
		err := addTypes(ht, w)
		if err != nil {
			return err
		}
	}
	fmt.Fprintf(w, funcTpl, funcName(r), jsonString(r.Path))
	return nil
}

func funcName(r *router.Route) string {
	name := r.GetName()
	if name == "" {
		name = strings.ReplaceAll(strings.TrimPrefix(r.Path, "/"), "/", ".")
	}
	out := ""
	upper := false
	for _, c := range name {
		if c == '.' {
			upper = true
			continue
		}
		if upper {
			upper = false
			out += strings.ToUpper(string(c))
			continue
		}
		out += string(c)
	}
	return out
}
