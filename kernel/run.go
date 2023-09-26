package kernel

import (
	"context"
	"fmt"
	"log"
	"net/http"
)

func (k *Kernel) Run(ctx context.Context) error {

	go k.RunListeners(ctx)
	go k.RunServices(ctx)

	return k.RunHttpServer(ctx)
}

func (k *Kernel) RunHttpServer(ctx context.Context) error {
	log.Printf("http://localhost:%d", k.port)

	handler := k.rootHandler()
	return http.ListenAndServe(fmt.Sprintf(":%d", k.port), handler)
}

func (k *Kernel) RunServices(ctx context.Context) error {
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
	return nil
}
