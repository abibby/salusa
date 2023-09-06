package dbtest

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/jmoiron/sqlx"
)

type Runner struct {
	open func() (*sqlx.DB, error)
	db   *sqlx.DB
}

func NewRunner(open func() (*sqlx.DB, error)) *Runner {
	return &Runner{
		open: open,
	}
}

type runner[T any] interface {
	Run(name string, cb func(t T)) bool
}

func (r *Runner) Run(t *testing.T, name string, cb func(t *testing.T, tx *sqlx.Tx)) bool {
	return run(r, t, name, cb)
}
func (r *Runner) RunBenchmark(t *testing.B, name string, cb func(t *testing.B, tx *sqlx.Tx)) bool {
	return run(r, t, name, cb)
}
func run[T testing.TB](r *Runner, t T, name string, cb func(t T, tx *sqlx.Tx)) bool {
	var err error
	if r.db == nil {
		r.db, err = r.open()
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to open database: %v", err)
			t.FailNow()
		}
	}
	tx, err := r.db.BeginTxx(context.Background(), nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to begin transaction: %v", err)
		t.FailNow()
	}

	var tAny any = t
	result := tAny.(runner[T]).Run(name, func(t T) {
		cb(t, tx)
	})

	err = tx.Rollback()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to rollback transaction: %v", err)
		t.FailNow()
	}
	return result
}
