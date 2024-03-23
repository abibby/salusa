//go:build !dev

package resources

import (
	"embed"

	"github.com/abibby/salusa/view"
)

//go:embed dist/*
var Content embed.FS

var View = view.Factory(Content, "**/*.html")
