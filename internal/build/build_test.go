package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIsUppercase(t *testing.T) {
	assert.True(t, IsUppercase('A'))
	assert.True(t, IsUppercase('Z'))
	assert.False(t, IsUppercase('a'))
	assert.False(t, IsUppercase('1'))
	assert.False(t, IsUppercase('_'))
}

func TestGetStruct(t *testing.T) {
	file := filepath.Join(t.TempDir(), "foo.go")
	require.NoError(t, os.WriteFile(file, []byte(`package foo

type Builder struct {
    query string
    wheres *Conditions
    havings *Conditions
    // comment
}

type Generic[T any] struct {
    x int
}
`), 0644))

	name, params, fields, err := GetStruct(file, 2)
	require.NoError(t, err)
	assert.Equal(t, "Builder", name)
	assert.Empty(t, params)
	assert.Equal(t, []string{"query"}, fields["string"])
	assert.Equal(t, []string{"wheres", "havings"}, fields["*Conditions"])

	name, params, _, err = GetStruct(file, 9)
	require.NoError(t, err)
	assert.Equal(t, "Generic", name)
	assert.Equal(t, "[T]", params)
}

func TestGetStructNoClosingBrace(t *testing.T) {
	file := filepath.Join(t.TempDir(), "foo.go")
	require.NoError(t, os.WriteFile(file, []byte(`package foo

type Foo struct {
`), 0644))

	_, _, _, err := GetStruct(file, 2)
	assert.Error(t, err)
}

func TestGetStructMissingFile(t *testing.T) {
	_, _, _, err := GetStruct(filepath.Join(t.TempDir(), "missing.go"), 0)
	assert.Error(t, err)
}

func TestReadSource(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, "a.go"), []byte("package a\n"), 0644))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "b.go"), []byte("package b\n"), 0644))

	src, err := ReadSource(dir)
	require.NoError(t, err)
	assert.Contains(t, src, "package a")
	assert.Contains(t, src, "package b")
}

func TestReadSourceMissingDir(t *testing.T) {
	_, err := ReadSource(filepath.Join(t.TempDir(), "missing"))
	assert.Error(t, err)
}

func TestReadSourceUnreadableFile(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, "a.go"), []byte("package a\n"), 0644))
	require.NoError(t, os.Mkdir(filepath.Join(dir, "sub"), 0644))

	_, err := ReadSource(dir)
	assert.Error(t, err)
}

func TestMainFunc(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, "foo.go"), []byte(`package foo

type Builder struct {
    query   string
    wheres  *Conditions
    havings *Conditions
}

// Where adds a condition
func (c *Conditions) Where(condition string) *Conditions { return c }

func (c *Conditions) And(condition string) *Conditions { return c }

func (c *Conditions) Or(condition string) *Conditions { return c }

func (c *Conditions) Having(condition string) *Conditions { return c }

func (c *Conditions) In(column string, values ...any) *Conditions { return c }

func (c *Conditions) unexported(condition string) *Conditions { return c }

func (c *Conditions) Clone() *Conditions { return c }

func (o *Other) Thing() *Other { return o }
`), 0644))

	t.Chdir(dir)
	t.Setenv("GOLINE", "2")
	t.Setenv("GOFILE", "foo.go")
	t.Setenv("GOPACKAGE", "foo")

	main()

	b, err := os.ReadFile(filepath.Join(dir, "generated_Builder.go"))
	require.NoError(t, err)
	src := string(b)

	assert.Contains(t, src, "func (b *Builder) Where(condition string) *Builder {")
	assert.Contains(t, src, "b.wheres = b.wheres.Where(condition)")
	assert.Contains(t, src, "func (b *Builder) Having(condition string) *Builder {\n\tb.havings = b.havings.Where(condition)")
	assert.Contains(t, src, "func (b *Builder) HavingAnd(condition string) *Builder {\n\tb.havings = b.havings.And(condition)")
	assert.Contains(t, src, "func (b *Builder) HavingOr(condition string) *Builder {\n\tb.havings = b.havings.Or(condition)")
	assert.Contains(t, src, "b.wheres = b.wheres.In(column, values...)")
	assert.NotContains(t, src, "unexported")
	assert.NotContains(t, src, "Clone")
	assert.NotContains(t, src, "Thing")
}

func TestMainFuncBadLine(t *testing.T) {
	t.Setenv("GOLINE", "not-a-number")
	assert.Panics(t, main)
}

func TestMainFuncMissingFile(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)
	t.Setenv("GOLINE", "0")
	t.Setenv("GOFILE", "nope.go")
	assert.Panics(t, main)
}

func TestMainFuncUnreadableSource(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, "foo.go"), []byte(`package foo

type Builder struct {
    wheres *Conditions
}

func (c *Conditions) Where(condition string) *Conditions { return c }
`), 0644))
	require.NoError(t, os.Mkdir(filepath.Join(dir, "generated_Builder.go"), 0644))

	t.Chdir(dir)
	t.Setenv("GOLINE", "2")
	t.Setenv("GOFILE", "foo.go")
	assert.Panics(t, main)
}

func TestMainFuncUnwritable(t *testing.T) {
	if os.Getuid() == 0 {
		t.Skip("running as root")
	}
	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, "foo.go"), []byte(`package foo

type Builder struct {
    wheres *Conditions
}

func (c *Conditions) Where(condition string) *Conditions { return c }
`), 0644))

	sub := filepath.Join(dir, "ro")
	require.NoError(t, os.Mkdir(sub, 0755))
	require.NoError(t, os.Chmod(sub, 0555))
	t.Cleanup(func() { _ = os.Chmod(sub, 0755) })
	t.Chdir(sub)
	t.Setenv("GOLINE", "2")
	t.Setenv("GOFILE", filepath.Join(dir, "foo.go"))
	assert.Panics(t, main)
}
