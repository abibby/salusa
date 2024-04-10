package generateclient

import (
	"io"
	"net/http"
	"reflect"
	"strings"

	"github.com/abibby/salusa/internal/helpers"
	"github.com/abibby/salusa/router"
)

type HandlerTypes interface {
	RequestType() reflect.Type
	ResponseType() reflect.Type
}

func addTypes(w io.Writer, r *router.Route, ht HandlerTypes) (string, error) {
	name := firstCap(funcName(r))
	_, err := w.Write([]byte("export type " + name + "Request = "))
	if err != nil {
		return "", err
	}
	err = toTsType(w, ht.RequestType(), nil)
	if err != nil {
		return "", err
	}
	_, err = w.Write([]byte("\nexport type " + name + "Response = "))
	if err != nil {
		return "", err
	}
	err = toTsType(w, ht.ResponseType(), map[reflect.Type]string{
		helpers.GetType[http.Handler](): "unknown",
	})
	if err != nil {
		return "", err
	}
	_, err = w.Write([]byte("\n"))
	if err != nil {
		return "", err
	}
	return name, nil
}
func firstCap(s string) string {
	return strings.ToUpper(s[0:1]) + s[1:]
}
