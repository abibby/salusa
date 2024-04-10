package generateclient

import (
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/abibby/salusa/router"
	"github.com/abibby/salusa/slices"
)

var paramRegex = regexp.MustCompile(`{(\w+)}`)

func GenerateClient(w io.Writer, r *router.Router) error {

	for _, route := range r.Routes() {
		err := routeTs(w, route)
		if err != nil {
			return err
		}
	}
	return nil
}

var funcTpl = `export async function %s(%s): Promise<%s> {
	const response = await fetcher(%s)
	if (response.status < 200 || response.status >= 300) {
		throw new Error("invalid status")
	}
	return await response.json()
}
`

func routeTs(w io.Writer, r *router.Route) error {
	var err error
	ht, hasTypes := r.GetHandler().(HandlerTypes)
	name := ""
	if hasTypes {
		name, err = addTypes(w, r, ht)
		if err != nil {
			return err
		}
	}
	params := urlParams(r)
	ret := "undefined"
	if hasTypes {
		if params != "" {
			params += ", "
		}
		params += "req: " + name + "Request"

		ret = name + "Response"
	}
	_, err = fmt.Fprintf(w, funcTpl, funcName(r), params, ret, buildPath(r))
	return err
}

func jsonString(s string) string {
	b, err := json.Marshal(s)
	if err != nil {
		panic(err)
	}
	return string(b)
}

func funcName(r *router.Route) string {
	name := r.GetName()
	if name == "" {
		name = strings.ToLower(r.Method) + "." + strings.ReplaceAll(strings.TrimPrefix(paramRegex.ReplaceAllString(r.Path, ""), "/"), "/", ".")
	}

	out := ""
	upper := false
	for _, c := range name {
		switch c {
		case '.', '_', ' ':
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

func urlParams(r *router.Route) string {
	matches := paramRegex.FindAllStringSubmatch(r.Path, -1)
	params := slices.Map(matches, func(p []string) string {
		return p[1] + ": string|number"
	})
	return strings.Join(params, ", ")
}
func buildPath(r *router.Route) string {
	return paramRegex.ReplaceAllStringFunc(jsonString(r.Path), func(s string) string {
		return `" + encodeURIComponent(` + s[1:len(s)-1] + `) + "`
	})
}
