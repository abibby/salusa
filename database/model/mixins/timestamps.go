package mixins

import (
	"context"
	"time"

	"github.com/abibby/salusa/database"
	"github.com/abibby/salusa/database/hooks"
)

type Timestamps struct {
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

var _ hooks.BeforeSaver = (*Timestamps)(nil)

// BeforeSave implements hooks.BeforeSaver.
func (t *Timestamps) BeforeSave(ctx context.Context, tx database.DB) error {
	now := time.Now()
	if (t.CreatedAt == time.Time{}) {
		t.CreatedAt = now
	}
	t.UpdatedAt = now
	return nil
}
