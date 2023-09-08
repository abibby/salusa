package pkg

import (
	"fmt"
	"regexp"
	"strings"
)

type Package struct {
	functionCalls []*FunctionCall
	imports       map[string]string
}

func New() *Package {
	return &Package{
		functionCalls: []*FunctionCall{},
		imports:       map[string]string{},
	}
}

func (p *Package) Add(function string, args ...any) *FunctionCall {
	// parts := strings.SplitN(function, ".", 2)
	i := strings.LastIndex(function, ".")
	pkgPath := function[0:i]
	functionName := function[i+1:]
	pkg := regexp.MustCompile(`[^\w]+`).ReplaceAllLiteralString(pkgPath, "_")
	fc := &FunctionCall{
		name:            pkg + "." + functionName,
		args:            args,
		hasError:        false,
		returnCount:     1,
		returnVariables: map[int]string{},
	}
	p.functionCalls = append(p.functionCalls, fc)
	p.imports[pkg] = pkgPath
	return fc
}

func (p *Package) Run() error {
	return nil
}

func (p *Package) ToGo() string {
	src := "package main\n\nimport (\n"
	for alias, path := range p.imports {
		if alias == path {
			src += fmt.Sprintf("\t%#v\n", path)
		} else {
			src += fmt.Sprintf("\t%s %#v\n", alias, path)
		}
	}
	src += ")\n\nfunc main() {\n"

	for _, functionCall := range p.functionCalls {
		src += functionCall.ToGo()
	}
	return src + "}\n"
}
