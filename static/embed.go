package static

import "embed"

//go:embed template/*
var Content embed.FS
