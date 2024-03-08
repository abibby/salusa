package kernel

import (
	"context"
	"net/http"
	"testing"

	"github.com/abibby/salusa/di"
	"github.com/jmoiron/sqlx"
)

func RunTest(t *testing.T, k *Kernel, name string, cb func(t *testing.T, h http.Handler, db *sqlx.DB)) {
	t.Run(name, func(t *testing.T) {
		ctx := context.Background()
		err := k.Bootstrap(ctx)
		if err != nil {
			t.Errorf("failed to bootstrap: %v", err)
			return
		}

		db, err := di.Resolve[*sqlx.DB](ctx)
		if err != nil {
			t.Errorf("no db in di: %w", err)
			return
		}

		cb(t, k.rootHandler(), db)
	})
}
