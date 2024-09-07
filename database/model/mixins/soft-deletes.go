package mixins

import (
	"time"

	"github.com/abibby/salusa/database"
	"github.com/abibby/salusa/database/builder"
)

type SoftDelete struct {
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
}

func (f *SoftDelete) Scopes() []*builder.Scope {
	return []*builder.Scope{
		SoftDeleteScope,
	}
}

var SoftDeleteScope = &builder.Scope{
	Name: "soft-deletes",
	Query: func(b *builder.Builder) *builder.Builder {
		return b.Where(b.GetTable()+".deleted_at", "=", nil)
	},
	Delete: func(next func(q *builder.Builder, tx database.DB) error) func(q *builder.Builder, tx database.DB) error {
		return func(q *builder.Builder, tx database.DB) error {
			return q.Update(tx, builder.Updates{
				"deleted_at": time.Now(),
			})
		}
	},
}
