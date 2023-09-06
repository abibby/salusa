package builder

var SoftDeletes = &Scope{
	Name: "soft-deletes",
	Apply: func(b *SubBuilder) *SubBuilder {
		return b.Where(b.GetTable()+".deleted_at", "=", nil)
	},
}
