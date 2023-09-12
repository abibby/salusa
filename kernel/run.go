package kernel

import (
	"fmt"
	"log"
	"net/http"
)

func (k *Kernel) Run() error {
	handler := k.rootHandler

	for _, m := range k.middleware {
		handler = m.Middleware(handler)
	}

	log.Printf("http://localhost:%d", k.port)
	return http.ListenAndServe(fmt.Sprintf(":%d", k.port), handler)
}
