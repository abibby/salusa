package kernel

import (
	"context"
	"fmt"
	"log"
	"net/http"
)

func (k *Kernel) Run(ctx context.Context) error {
	handler := k.rootHandler

	for _, m := range k.middleware {
		handler = m.Middleware(handler)
	}
	for _, s := range k.services {
		go func(s Service) {
			for {
				err := s.Run(ctx, k)
				if err != nil {
					log.Print(err)
				}
				if !s.Restart() {
					return
				}
			}
		}(s)
	}

	go k.runListeners()

	log.Printf("http://localhost:%d", k.port)

	return http.ListenAndServe(fmt.Sprintf(":%d", k.port), handler)
}
