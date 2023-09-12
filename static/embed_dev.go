//go:build dev

package static

import (
	"io/fs"
	"os"
)

var Content fs.FS = os.DirFS("template")
