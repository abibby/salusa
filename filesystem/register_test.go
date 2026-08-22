package filesystem_test

import (
	"io/fs"
	"os"
	"path/filepath"
	"testing"
	"testing/fstest"

	"github.com/abibby/salusa/di"
	"github.com/abibby/salusa/filesystem"
	"github.com/abibby/salusa/salusaconfig"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewLocalFS(t *testing.T) {
	dir := t.TempDir()
	err := os.WriteFile(filepath.Join(dir, "file.txt"), []byte("hello"), 0644)
	require.NoError(t, err)

	l := filesystem.NewLocalFS(dir)
	require.NotNil(t, l)
	assert.Equal(t, dir, l.Root)

	b, err := fs.ReadFile(l.FS(), "file.txt")
	require.NoError(t, err)
	assert.Equal(t, "hello", string(b))
}

func TestLocalFS_FS(t *testing.T) {
	l := filesystem.LocalFS{Root: "testdata"}
	require.NotNil(t, l.FS())
}

type testConfig struct {
	salusaconfig.Config
	fsys fs.FS
}

func (t *testConfig) FSConfig() filesystem.Config {
	return &localConfig{fsys: t.fsys}
}

func (t *testConfig) FS() fs.FS {
	return t.fsys
}

type localConfig struct {
	fsys fs.FS
}

func (l *localConfig) FS() fs.FS {
	return l.fsys
}

func TestRegister(t *testing.T) {
	fsys := fstest.MapFS{
		"file.txt": &fstest.MapFile{Data: []byte("hello")},
	}

	ctx := di.TestDependencyProviderContext()
	di.RegisterSingleton(ctx, func() salusaconfig.Config {
		return &testConfig{fsys: fsys}
	})

	filesystem.Register(ctx)

	f, err := di.Resolve[fs.FS](ctx)
	require.NoError(t, err)

	b, err := fs.ReadFile(f, "file.txt")
	require.NoError(t, err)
	assert.Equal(t, "hello", string(b))
}

type nonFSConfig struct {
	salusaconfig.Config
}

func TestRegister_notFSConfiger(t *testing.T) {
	ctx := di.TestDependencyProviderContext()
	di.RegisterSingleton(ctx, func() salusaconfig.Config {
		return &nonFSConfig{}
	})

	filesystem.Register(ctx)

	_, err := di.Resolve[fs.FS](ctx)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "config not instance of email.FSConfiger")
}
