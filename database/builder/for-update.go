package builder

import "abibby.com/salusa/database/dialects"

func (b *Builder) ForUpdate() *Builder {
	b.query.ForUpdate = dialects.ForUpdateDefault
	return b
}
func (b *Builder) ForUpdateSkipLocked() *Builder {
	b.query.ForUpdate = dialects.ForUpdateSkipLocked
	return b
}
