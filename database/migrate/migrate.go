package migrate

import (
	"fmt"
	"go/format"

	"abibby.com/salusa/database/schema"
	"golang.org/x/tools/imports"
)

type Migration struct {
	Name string
	Up   schema.Runner
	Down schema.Runner
}

func SrcFile(migrationName, packageName string, up, down fmt.GoStringer) (string, error) {
	outFile := "migration.go"
	initSrc := `package %s
	
import (
	"abibby.com/salusa/database"
	"abibby.com/salusa/database/migrate"
	"abibby.com/salusa/database/builder"
	"abibby.com/salusa/database/schema"
)

func init() {
	migrations.Add(&migrate.Migration{
		Name: %#v,
		Up: %s,
		Down: %s,
	})
}`

	src := []byte(fmt.Sprintf(initSrc, packageName, migrationName, up.GoString(), down.GoString()))
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
