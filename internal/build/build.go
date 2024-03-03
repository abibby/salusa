package main

import (
	"fmt"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/tools/imports"
)

func main() {

	line, err := strconv.Atoi(os.Getenv("GOLINE"))
	if err != nil {
		panic(err)
	}
	file := os.Getenv("GOFILE")
	pkg := os.Getenv("GOPACKAGE")

	structName, structParams, structFields, err := GetStruct(file, line)
	if err != nil {
		panic(err)
	}

	goSrc, err := ReadSource(".")
	if err != nil {
		panic(err)
	}
	// matches := regexp.MustCompile(`\nfunc \((\w+ +)?([^)]+)\) ([\w)]+)\((.*)\) (\*?\w+(?:\[.+\])?) {`).FindAllStringSubmatch(goSrc, -1)
	matches := regexp.MustCompile(`((?:\/\/[^\n]+\n)*)func \((\w+ +)?([^)]+)\) ([\w)]+)\((.*)\) (\*?\w+(?:\[.+\])?) {`).FindAllStringSubmatch(goSrc, -1)

	fmt.Printf("Generating code for %s in %s:%d\n", structName, file, line)

	src := "package " + pkg + "\n\n"
	for _, match := range matches {
		comment := match[1]
		fieldType := match[3]
		methodName := match[4]
		params := match[5]
		returnType := match[6]

		if methodName == "Clone" || methodName == "ToSQL" || !IsUppercase(methodName[0]) || returnType != fieldType {
			continue
		}

		if fieldNames, ok := structFields[fieldType]; ok {
			for _, fieldName := range fieldNames {
				originalMethodName := methodName
				if fieldName == "havings" {
					originalMethodName := methodName
					switch methodName {
					case "And":
						methodName = "HavingAnd"
					case "Or":
						methodName = "HavingOr"
					default:
						methodName = strings.ReplaceAll(methodName, "Where", "Having")
						comment = strings.ReplaceAll(comment, "where", "having")
					}
					comment = strings.ReplaceAll(comment, originalMethodName, methodName)
				}
				args := ""
				for i, p := range strings.Split(params, ",") {
					if i != 0 {
						args += ", "
					}
					parts := strings.Split(strings.TrimSpace(p), " ")
					args += parts[0]
					if len(parts) > 1 && strings.HasPrefix(parts[1], "...") {
						args += "..."
					}
				}
				src += fmt.Sprintf(
					"%sfunc (b *%s) %s(%s) *%s {\n"+
						"\tb = b.Clone()\n"+
						"\tb.%s = b.%s.%s(%s)\n"+
						"\treturn b\n"+
						"}\n",
					comment,
					structName+structParams,
					methodName,
					params,
					structName+structParams,
					fieldName,
					fieldName,
					originalMethodName,
					args,
				)
			}
		}
	}

	outFile := fmt.Sprintf("generated_%s.go", structName)

	b, err := imports.Process(outFile, []byte(src), nil)
	if err != nil {
		panic(err)
	}

	err = os.WriteFile(outFile, b, 0644)
	if err != nil {
		panic(err)
	}

}

func IsUppercase(r byte) bool {
	return r >= 'A' && r <= 'Z'
}

func GetStruct(file string, line int) (string, string, map[string][]string, error) {
	m := map[string][]string{}

	b, err := os.ReadFile(file)
	if err != nil {
		return "", "", nil, err
	}
	src := string(b)

	lines := strings.Split(src, "\n")

	matches := regexp.MustCompile(`type (\w*)(\[.*\])? struct {`).FindStringSubmatch(lines[line])
	structName := matches[1]
	params := matches[2]
	if params != "" {
		params = "[T]"
	}

	for _, l := range lines[line+1:] {
		if l == "}" {
			return structName, params, m, nil
		}
		parts := strings.Fields(strings.TrimSpace(l))
		if len(parts) < 2 {
			continue
		}
		a, ok := m[parts[1]]
		if !ok {
			a = make([]string, 0, 1)
		}
		a = append(a, parts[0])
		m[parts[1]] = a
	}

	return "", "", nil, fmt.Errorf("no struct found in %s:%d", file, line)
}

func ReadSource(root string) (string, error) {
	dir, err := os.ReadDir(root)
	if err != nil {
		return "", err
	}
	src := ""
	for _, f := range dir {
		b, err := os.ReadFile(path.Join(root, f.Name()))
		if err != nil {
			return "", err
		}
		src += string(b)
	}
	return src, nil
}
