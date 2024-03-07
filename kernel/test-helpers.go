package kernel

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"testing"
)

func RunTest(t *testing.T, k *Kernel, name string, cb func(t *testing.T, h http.Handler)) {
	t.Run(name, func(t *testing.T) {
		ctx := context.Background()
		err := k.Bootstrap(ctx)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to bootstrap: %v", err)
			t.Fail()
			return
		}

		cb(t, k.rootHandler())
	})
}
