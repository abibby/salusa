package util

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfig_prep_withImport(t *testing.T) {
	c := Config{
		Model: &Package{
			Dir:    "app/models",
			Pkg:    "models",
			Import: "custom/import",
		},
		Migration: &Package{Dir: "migrations"},
	}
	c.prep("/root", &PackageInfo{RootPackage: "mod", SpicePackage: "mod"})

	require.NotNil(t, c.Model)
	assert.Equal(t, "custom/import", c.Model.Import)
	require.NotNil(t, c.Migration)
	assert.Equal(t, "mod/migrations", c.Migration.Import)
}

func TestConfig_prepWithSubdirRoot(t *testing.T) {
	c := Config{
		Model:     &Package{Dir: "app/models"},
		Migration: &Package{Dir: "migrations"},
	}
	c.prep("/root/app", &PackageInfo{RootPackage: "mod", SpicePackage: "mod/app"})

	require.NotNil(t, c.Model)
	assert.Equal(t, filepath.Join("/root/app", "app", "models"), c.Model.Dir)
	assert.True(t, strings.HasPrefix(c.Model.Import, "mod/app/"))
}
