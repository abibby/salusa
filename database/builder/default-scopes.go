package builder

var SoftDeletes = &Scope{
	Name: "soft-deletes",
	Apply: func(b *Builder) *Builder {
		return b.Where(b.GetTable()+".deleted_at", "=", nil)
	},
}
