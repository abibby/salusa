package generateclient

import (
	"fmt"
	"io"
)

type HandlerTypes interface {
	GetRequest() any
	GetResponse() any
}

func addTypes(ht HandlerTypes, w io.Writer) error {
	fmt.Fprintf(w, "%#v\n", ht.GetRequest())
	return nil
}
