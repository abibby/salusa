package util

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/abibby/salusa/slices"
	"golang.org/x/mod/modfile"
)

var ErrNoPackage = fmt.Errorf("not in a go package")

func PackageRoot(from string) (string, string, error) {
	cwd, err := filepath.Abs(from)
	if err != nil {
		return "", "", err
	}
	current := cwd
	name := ""
	relPackage := ""
	for {
		files, err := os.ReadDir(current)
		if err != nil {
			return "", "", err
		}

		_, ok := slices.Find(files, func(f fs.DirEntry) bool {
			return f.Name() == "go.mod"
		})
		if ok {
			return current, relPackage, nil
		}

		if current == "/" {
			return "", "", ErrNoPackage
		}
		current = strings.TrimSuffix(current, "/")
		current, name = path.Split(current)
		relPackage = path.Join(name, relPackage)
	}
}

type PackageInfo struct {
	PackageRoot    string
	RootPackage    string
	CurrentPackage string
}

func PkgInfo(from string) (*PackageInfo, error) {
	root, relPackage, err := PackageRoot(from)
	if err != nil {
		return nil, err
	}

	modFile := path.Join(root, "go.mod")
	b, err := os.ReadFile(modFile)
	if err != nil {
		return nil, err
	}

	m, err := modfile.ParseLax(modFile, b, nil)
	if err != nil {
		return nil, err
	}

	pkg := m.Module.Syntax.Token[1]
	return &PackageInfo{
		PackageRoot:    root,
		RootPackage:    pkg,
		CurrentPackage: path.Join(pkg, relPackage),
	}, nil
}
