package migrate

import (
	"fmt"
	"go/format"

	"github.com/abibby/salusa/database/schema"
	"golang.org/x/tools/imports"
)

type Migration struct {
	Name string
	Up   schema.Runner
	Down schema.Runner
}

func SrcFile(migrationName, packageName string, up, down ToGoer) (string, error) {
	outFile := "test.go"
	initSrc := `package %s
	
import (
	"github.com/abibby/salusa/database/migrate"
	"github.com/abibby/salusa/database/builder"
	"github.com/abibby/salusa/database/schema"
)

func init() {
	migrations.Add(&migrate.Migration{
		Name: %#v,
		Up: %s,
		Down: %s,
	})
}`

	src := []byte(fmt.Sprintf(initSrc, packageName, migrationName, up.ToGo(), down.ToGo()))
	// fmt.Printf("%s\n", src)
	src, err := imports.Process(outFile, src, nil)
	if err != nil {
		return "", err
	}

	src, err = format.Source(src)
	if err != nil {
		return "", err
	}
	return string(src), nil
}
