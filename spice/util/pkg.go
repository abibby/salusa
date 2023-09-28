package util

import (
	"os"
	"path"
	"strings"

	"golang.org/x/mod/modfile"
)

type PackageInfo struct {
	RootPackage  string
	SpicePackage string
}

func PkgInfo(goRoot, spiceRoot string) (*PackageInfo, error) {
	modFile := path.Join(goRoot, "go.mod")
	b, err := os.ReadFile(modFile)
	if err != nil {
		return nil, err
	}

	m, err := modfile.ParseLax(modFile, b, nil)
	if err != nil {
		return nil, err
	}

	pkg := m.Module.Syntax.Token[1]

	spicePkg := pkg
	p, ok := strings.CutPrefix(spiceRoot, goRoot)
	if ok {
		spicePkg = path.Join(pkg, p)
	}

	return &PackageInfo{
		RootPackage:  pkg,
		SpicePackage: spicePkg,
	}, nil
}
